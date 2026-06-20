# CEL Plugin

Transforms messages with CEL expressions.

Config:

```yaml
components:
  - id: "imap_to_telegram_cel_transformer"
    type: processor
    plugin:
      path: ./build/bin/cel-plugin.so
    config:
      expr: |
        {
          "chat_id": 1234,
          "content": self.?Envelope.?Subject.orValue(null)
        }
```

