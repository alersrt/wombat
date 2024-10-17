package main

import (
	"log"
	"wombat/internal/config"
	"wombat/pkg/daemon"
)

func main() {
	daemon.Create(&config.Config{}, func() {
		log.Println("asdasd")
	})
}
