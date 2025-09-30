module telegram

require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1

replace (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/OvyFlash/telegram-bot-api v0.0.0-20250903213241-2ddbaeebe9a5
	github.com/wombat => ../../
)
