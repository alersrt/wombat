module wombat

go 1.25.1

require (
	github.com/andygrunwald/go-jira v1.16.0
	github.com/emersion/go-imap/v2 v2.0.0-beta.6
	github.com/emersion/go-message v0.18.2
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	gopkg.in/yaml.v2 v2.4.0
	mvdan.cc/sh/v3 v3.12.0
)

require (
	github.com/emersion/go-sasl v0.0.0-20241020182733-b788ff22d5a6 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/trivago/tgo v1.0.7 // indirect
	golang.org/x/text v0.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/OvyFlash/telegram-bot-api v0.0.0-20250723080846-4b7dce3e82ce
