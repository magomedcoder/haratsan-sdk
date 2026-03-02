#!/usr/bin/env python3
from __future__ import annotations

import logging
import os
import signal
import sys
import threading
from datetime import datetime

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from haratsan_sdk import Client, run_polling, Update, build_reply_markup

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


def main() -> None:
    client = Client(addr="хост", token="токен")
    stop = threading.Event()
    signal.signal(signal.SIGINT, lambda *_: stop.set())

    def handler(update: Update) -> None:
        content = update.content.strip()
        if content == "/start":
            client.send_message(update.from_user_id, "Привет, я пример бота")

        elif content == "/time":
            reply = f"Сейчас время {datetime.now()}"
            message_id = client.send_message(update.from_user_id, reply)
            logger.info("Ответ пользователю %s (message_id=%s): %r", update.from_user_id, message_id, reply)

        elif content == "/vote":
            markup = build_reply_markup([
                [
                    ("Да", "vote_yes"), 
                    ("Нет", "vote_no"),
                ],
            ])
            client.send_message(update.from_user_id, "Голосуйте:", reply_markup=markup)
            logger.info("Отправлено сообщение с кнопками пользователю %s", update.from_user_id)

        elif content == "/menu":
            markup = build_reply_markup([
                [
                    ("Справка", "help"),
                    ("Время", "time"),
                ],
                [
                    ("Отмена", "cancel")
                ],
            ])
            client.send_message(update.from_user_id, "Выберите действие:", reply_markup=markup)

        else:
            try:
                message_id = client.send_message(update.from_user_id, content)
                logger.info("Ответ пользователю %s (message_id=%s): %r", update.from_user_id, message_id, content)
            except Exception as e:
                logger.exception("SendMessage: %s", e)

    run_polling(client, handler, stop_event=stop)


if __name__ == "__main__":
    main()
