package httpserv

import (
	"MinecraftServerLogger/config"
	"MinecraftServerLogger/logger"
	"MinecraftServerLogger/ws"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

type LoggerHTTPServer struct {
	conf   *config.RootConfig
	logger *logger.Logger
}

func NewLoggerHTTPServer(conf *config.RootConfig, logger *logger.Logger) *LoggerHTTPServer {
	return &LoggerHTTPServer{
		conf:   conf,
		logger: logger,
	}
}

func (l *LoggerHTTPServer) StartServer() error {
	router := mux.NewRouter()
	router.HandleFunc("/ws/log", l.websocketHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	addr := fmt.Sprintf("%s:%d", l.conf.HTTPConfig.ListenAddr, l.conf.HTTPConfig.ListenPort)

	if l.conf.HTTPConfig.UseSSL {
		return http.ListenAndServeTLS(addr,
			l.conf.HTTPConfig.CertFilePath,
			l.conf.HTTPConfig.KeyFilePath,
			router,
		)
	} else {
		return http.ListenAndServe(addr, router)
	}
}

func (l *LoggerHTTPServer) websocketHandler(w http.ResponseWriter, r *http.Request) {
	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade: ", err)
		return
	}

	wsReceiver := ws.NewWebsocketLoggerReceiver(l.logger, wsconn)
	l.logger.RegisterNewReceiver(wsReceiver)
}
