#!/usr/bin/env python3
import logging
import os
import signal
import sys
import threading

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from haratsan_sdk import Client, run_polling, Update

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def main() -> None:
    client = Client(addr="хост", token="токен")
    stop = threading.Event()
    signal.signal(signal.SIGINT, lambda *_: stop.set())

    def handler(update: Update) -> None:
        content = update.content.strip()
        if content == "/start":
            reply = "Привет, я пример бота"
        elif content == "/time":
            from datetime import datetime
            reply = f"Сейчас время {datetime.now()}"
        else:
            reply = content

        try:
            message_id = client.send_message(update.from_user_id, reply)
            logger.info(
                "Ответ пользователю %s (message_id=%s, sent_id=%s): %r",
                update.from_user_id,
                update.message_id,
                message_id,
                reply,
            )
        except Exception as e:
            logger.exception("SendMessage: %s", e)

    run_polling(client, handler, stop_event=stop)


if __name__ == "__main__":
    main()
