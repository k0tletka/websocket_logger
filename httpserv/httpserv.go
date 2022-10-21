package httpserv

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/k0tletka/websocket_logger/config"
	"github.com/k0tletka/websocket_logger/logger"
	"github.com/k0tletka/websocket_logger/ws"
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

	router.Use(l.basicAuth)

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

func (l *LoggerHTTPServer) basicAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()

		if !ok || !l.checkBasicAuthCredentials(user, password) {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"Please, provide valid username and password\"")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (l *LoggerHTTPServer) checkBasicAuthCredentials(user, password string) bool {
	for _, account := range l.conf.HTTPConfig.BasicAuthUsers {
		passwordBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(password)))
		base64.StdEncoding.Encode(passwordBase64, []byte(password))

		if account.Name == user && account.Base64Hash == string(passwordBase64) {
			return true
		}
	}

	return false
}
