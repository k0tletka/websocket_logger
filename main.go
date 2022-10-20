package main

import (
	"MinecraftServerLogger/config"
	"MinecraftServerLogger/logger"
	"log"
)

func main() {
	cfg, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	tailLogger := logger.NewLogger(cfg)
	if err := tailLogger.Start(); err != nil {
		log.Fatalln(err)
	}
}
