Wombat
======

Why Wombat? Working bot for messages, w-m-bot.

Description
-----------

DevOps related app needed to react on email's alerts.

Main features:
- Read alerted emails
- Send alerts to Telegram
- Create tasks for infrastructure engineers in Jira and update sent alerts with related Jira's links

Usage
------

It's necessary to specify path to custom [config](examples/config.yaml): `--config=/path/to/config`. There is possibility to specify environment variables in config also.

### Default environment variables (can be changed)

| Env                        | Description | Default                            |
|:---------------------------|:------------|:-----------------------------------|
| `WOMBAT_JIRA_URL`          |             |                                    |
| `WOMBAT_TELEGRAM_TOKEN`    |             |                                    |

[extra]: https://github.com/golang-standards/project-layout
