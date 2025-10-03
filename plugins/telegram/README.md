# Telegram Plugin

Sends messages to specified Telegram's chat.

Config:

```yaml
consumers:
    - name: tg_consumer
      plugin: telegram
      conf:
          token: ${WOMBAT_TELEGRAM_TOKEN:-example}
```

Consumed format:

```json
{
    "chat_id": 9999999999,
    "content": "<message content>"
}
```
