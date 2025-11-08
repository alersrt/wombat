# IMAP Plugin

Gets messages from IMAP.

Config:

```yaml
producers:
    - name: imap_producer
      plugin: imap
      conf:
          url: ${WOMBAT_IMAP_ADDRESS:-mail.dev:993}
          username: ${WOMBAT_IMAP_USERNAME:-example@mail.dev}
          password: ${WOMBAT_IMAP_PASSWORD:-password}
          mailbox: ${WOMBAT_IMAP_MAILBOX:-Inbox}
          idleTimeout: ${WOMBAT_IMAP_IDLE_TIMEOUT:-5000}
          verbose: ${WOMBAT_IMAP_VERBOSE:-true}
```

Produced format:

```json
{
    "text": "<message content>",
    "envelope": {
        "date": "<RFC3339>",
        "subject": "string",
        "from": {
            "name": "string",
            "mailbox": "string",
            "host": "string"
        },
        "sender": {
            "name": "string",
            "mailbox": "string",
            "host": "string"
        },
        "reply_to": {
            "name": "string",
            "mailbox": "string",
            "host": "string"
        },
        "to": {
            "name": "string",
            "mailbox": "string",
            "host": "string"
        },
        "cc": {
            "name": "string",
            "mailbox": "string",
            "host": "string"
        },
        "bcc": {
            "name": "string",
            "mailbox": "string",
            "host": "string"
        },
        "in_reply_to": ["string"],
        "message_id": "string"
    }
}
```
