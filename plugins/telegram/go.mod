module telegram

go 1.25.1

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/wombat v0.0.0-00010101000000-000000000000
)

replace (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/OvyFlash/telegram-bot-api v0.0.0-20250903213241-2ddbaeebe9a5
	github.com/wombat => ../../
)
