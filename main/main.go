package main

import (
	"log"
	"wombat/internal/config"
	"wombat/pkg/daemon"
)

func main() {
	conf := &config.Config{}
	daemon.Create(conf, func() {
		log.Println("asdasd")
	})
}
