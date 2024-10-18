package main

import (
	"log"
	"time"
	"wombat/internal/config"
	"wombat/pkg/daemon"
)

func main() {
	conf := &config.Config{}
	daemon.Create(conf, func() {
		time.Sleep(1 * time.Second)
		log.Println("<<-->>")
	})
}
