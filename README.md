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

Auth
----

Right now it supports only really simple auth: bind to telegram's `user_id`.

Usage
------

It's necessary to specify path to custom [config](./cmd/config.yaml): `--config=/path/to/config`. There is possibility to specify environment variables in config also.

### Default environment variables (can be changed)

| Env                        | Description | Default      |
|:---------------------------|:------------|:-------------|
| `WOMBAT_TAG_PATTERN`       |             | `(TEST-\d+)` |
| `WOMBAT_JIRA_URL`          |             |              |
| `WOMBAT_TELEGRAM_TOKEN`    |             |              |
| `WOMBAT_CIPHER_KEY`        |             | `0==H`       |
| `WOMBAT_POSTGRES_USERNAME` |             | `postgres`   |
| `WOMBAT_POSTGRES_PASSWORD` |             | `postgres`   |
| `WOMBAT_POSTGRES_HOST`     |             | `localhost`  |
| `WOMBAT_POSTGRES_PORT`     |             | `5432`       |
| `WOMBAT_POSTGRES_DATABASE` |             | `postgres`   |
| `WOMBAT_POSTGRES_SSLMODE`  |             | `disable`    |

[extra]: https://github.com/golang-standards/project-layout
