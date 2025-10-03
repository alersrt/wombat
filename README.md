# Wombat

Why Wombat? Working bot for messages, w-m-bot.

## Description

Started as simple bot for messaging from telegram to jira, mutated to DevOps related app needed to react on email's alerts, and finally finished as extendable ETL platform.

Main features:

-   Plugins written on Go
-   Ready to Go plugins for IMAP, Telegram and Jira

## Usage

It's necessary to specify path to custom [config](examples/config.yaml): `--config=/path/to/config`. There is possibility to specify environment variables in config also.

## Plugins

Plugin should implement `Plugin` interface from [`pkg/interfaces.go`](pkg/interfaces.go) **and** `Producer` or `Consumer` or both.
Plugin should also has special function for exporting:

```go
type MyAwesomePlugin struct {

}

func Export() pkg.Plugin {
	return &MyAwesomePlugin{}
}
```

Build flags:

```shell
go build --trimpath --ldflags="-w -s" --buildmode=plugin -o <plugin name>.so
```

[cel-spec]: https://github.com/google/cel-spec
[cel-k8s]: https://kubernetes.io/docs/reference/using-api/cel/
