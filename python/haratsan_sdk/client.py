from __future__ import annotations

import threading
import time
from dataclasses import dataclass
from typing import Callable, Optional

import grpc

from .gen_pb import bot_api_pb2, bot_api_pb2_grpc

DEFAULT_POLL_INTERVAL = 2.0
DEFAULT_UPDATES_LIMIT = 50


@dataclass
class Update:
    message_id: int
    from_user_id: int
    content: str
    created_at: int


class _BotTokenInterceptor(grpc.UnaryUnaryClientInterceptor):
    def __init__(self, token: str):
        self._token = token

    def intercept_unary_unary(self, continuation, client_call_details, request):
        metadata = list(client_call_details.metadata or [])
        metadata.append(("bot-token", self._token))
        new_details = client_call_details._replace(metadata=metadata)
        return continuation(new_details, request)


class Client:
    def __init__(self, addr: str, token: str):
        self._token = token
        self._channel = grpc.insecure_channel(addr)
        interceptors = [_BotTokenInterceptor(token)]
        self._intercepted_channel = grpc.intercept_channel(self._channel, *interceptors)
        self._stub = bot_api_pb2_grpc.BotApiServiceStub(self._intercepted_channel)

    def close(self) -> None:
        self._channel.close()

    def __enter__(self) -> "Client":
        return self

    def __exit__(self, *args: object) -> None:
        self.close()

    def get_updates(
        self,
        offset: int,
        limit: int,
        timeout: Optional[float] = None,
    ) -> list[Update]:
        req = bot_api_pb2.GetUpdatesRequest(offset=offset, limit=limit)
        kwargs = {}
        if timeout is not None:
            kwargs["timeout"] = timeout
        resp = self._stub.GetUpdates(req, **kwargs)
        return [
            Update(
                message_id=u.message_id,
                from_user_id=u.from_user_id,
                content=u.content,
                created_at=u.created_at,
            )
            for u in resp.updates
        ]

    def send_message(
        self,
        to_user_id: int,
        content: str,
        timeout: Optional[float] = None,
    ) -> int:
        req = bot_api_pb2.SendMessageRequest(to_user_id=to_user_id, content=content)
        kwargs = {}
        if timeout is not None:
            kwargs["timeout"] = timeout
        resp = self._stub.SendMessage(req, **kwargs)
        return resp.message_id


UpdateHandler = Callable[[Update], None]


def run_polling(
    client: Client,
    handler: UpdateHandler,
    *,
    poll_interval: float = DEFAULT_POLL_INTERVAL,
    updates_limit: int = DEFAULT_UPDATES_LIMIT,
    stop_event: Optional[threading.Event] = None,
) -> None:
    if poll_interval <= 0:
        poll_interval = DEFAULT_POLL_INTERVAL
    if updates_limit <= 0:
        updates_limit = DEFAULT_UPDATES_LIMIT

    offset = 0
    while True:
        if stop_event is not None and stop_event.is_set():
            return
        try:
            updates = client.get_updates(offset, updates_limit)
        except grpc.RpcError:
            if stop_event is not None:
                if stop_event.wait(timeout=poll_interval):
                    return
            else:
                time.sleep(poll_interval)
            continue

        for u in updates:
            if u.message_id <= 0:
                continue
            offset = u.message_id
            try:
                handler(u)
            except Exception:
                pass

        if stop_event is not None:
            if stop_event.wait(timeout=poll_interval):
                return
        else:
            time.sleep(poll_interval)
