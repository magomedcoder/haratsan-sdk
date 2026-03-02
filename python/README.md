# Haratsan SDK Python

Клиент для бот-API Haratsan на Python.

Требуется Python 3.9+.

Установка:

```bash
pip install git+https://github.com/magomedcoder/haratsan-sdk.git#subdirectory=python
```

В коде:

```python
from haratsan_sdk import Client, run_polling, Update

client = Client(addr="хост", token="токен")

def handler(update: Update) -> None:
    client.send_message(update.from_user_id, "Ответ")

run_polling(client, handler)
```

### Inline-клавиатура (ReplyMarkup)

```python
from haratsan_sdk import Client, build_reply_markup

client = Client(addr="хост", token="токен")

markup = build_reply_markup([
    [
        ("Да", "vote_yes"), 
        ("Нет", "vote_no"),
    ],
])
client.send_message(user_id, "Голосуйте:", reply_markup=markup)

markup = build_reply_markup([
    [
        ("Справка", "help"), 
        ("Время", "time")
    ],
    [
        ("Отмена", "cancel")
    ],
])
client.send_message(user_id, "Выберите:", reply_markup=markup)
```

**Пример использования:** [`example/main.py`](example/main.py)