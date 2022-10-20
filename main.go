package main

import (
	"MinecraftServerLogger/config"
	"MinecraftServerLogger/httpserv"
	"MinecraftServerLogger/logger"
	"log"
)

func main() {
	cfg, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("GetConfiguration: ", err)
	}

	tailLogger := logger.NewLogger(cfg)
	if err := tailLogger.Start(); err != nil {
		log.Fatalln("Logger: ", err)
	}

	server := httpserv.NewLoggerHTTPServer(cfg, tailLogger)
	if err := server.StartServer(); err != nil {
		log.Fatalln("HTTPServer: ", err)
	}
}
