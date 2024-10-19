module wombat

go 1.20

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/OvyFlash/telegram-bot-api/v5 v5.0.0-20241013102643-36756d99d4ae
