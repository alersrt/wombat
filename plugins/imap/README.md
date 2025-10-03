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
    "Text": "<message content>",
    "Envelope": {
        "Date": "<RFC3339>",
        "Subject": "string",
        "From": {
            "Name": "string",
            "Mailbox": "string",
            "Host": "string"
        },
        "Sender": {
            "Name": "string",
            "Mailbox": "string",
            "Host": "string"
        },
        "ReplyTo": {
            "Name": "string",
            "Mailbox": "string",
            "Host": "string"
        },
        "To": {
            "Name": "string",
            "Mailbox": "string",
            "Host": "string"
        },
        "Cc": {
            "Name": "string",
            "Mailbox": "string",
            "Host": "string"
        },
        "Bcc": {
            "Name": "string",
            "Mailbox": "string",
            "Host": "string"
        },
        "InReplyTo": ["string"],
        "MessageID": "string"
    }
}
```
