Wombat
======

Why Wombat? Working bot for messengers, w-m-bot.

Description
-----------

The main idea of this bot is providing possibility to add comments/texts to Jira tasks directly from messenger.

The main idea:

1. User sends message directly to bot.
2. Bot get messages from API.
3. Parses message to match issues' numbers pattern.
4. Creates comment in related Jira task or edits existed one.

```mermaid
flowchart LR
    user([User])
    kafka([Kafka])
    tapi([Telegram API])
    router[Router]
    config[Config]
    cc[Confluence Connector]
    jc[Jira Connector]
    
    
    user --> tapi
    tapi --> router
    config --> router
    router --> kafka
    kafka --> cc
    kafka --> jc
```

Auth
----

Right now it supports only really simple auth: small table with list of usernames and flags of allowing.

Usage
------

It's necessary to specify path to custom [config](./cmd/config.yaml): `--config=/path/to/config`. There is possibility to specify environment variables in config also.

### Default environment variables (can be changed)

| Env                 | Description | Default           |
|:--------------------|:------------|:------------------|
| `TAG_PATTERN`       |             | `(TEST-\d+)`      |
| `JIRA_URL`          |             |                   |
| `JIRA_USERNAME`     |             |                   |
| `JIRA_TOKEN`        |             |                   |
| `TELEGRAM_TOKEN`    |             |                   |
| `POSTGRES_USERNAME` |             | `wombat_rw`       |
| `POSTGRES_PASSWORD` |             | `wombat_rw`       |
| `POSTGRES_HOST`     |             | `localhost`       |
| `POSTGRES_PORT`     |             | `5432`            |
| `POSTGRES_DATABASE` |             | `wombatdb`        |
| `POSTGRES_SSLMODE`  |             | `disable`         |
| `KAFKA_GROUP_ID`    |             | `wombat`          |
| `KAFKA_BOOTSTRAP`   |             | `localhost:9092`  |
| `KAFKA_TOPIC`       |             | `wombat.response` |

[extra]: https://github.com/golang-standards/project-layout
