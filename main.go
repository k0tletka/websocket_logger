package main

import (
	"github.com/k0tletka/websocket_logger/config"
	"github.com/k0tletka/websocket_logger/httpserv"
	"github.com/k0tletka/websocket_logger/logger"
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
