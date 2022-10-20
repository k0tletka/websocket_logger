package ws

import (
	"MinecraftServerLogger/logger"
	"github.com/gorilla/websocket"
	"log"
)

type WebsocketResponse struct {
	History []string `json:"history,omitempty"`
	Message string   `json:"message,omitempty"`
}

type WebsocketLoggerReceiver struct {
	logger *logger.Logger
	conn   *websocket.Conn
}

func NewWebsocketLoggerReceiver(logger *logger.Logger, conn *websocket.Conn) *WebsocketLoggerReceiver {
	return &WebsocketLoggerReceiver{
		logger: logger,
		conn:   conn,
	}
}

func (w *WebsocketLoggerReceiver) ReceiveMessage(message string) {
	resp := &WebsocketResponse{Message: message}

	if err := w.conn.WriteJSON(resp); err != nil {
		log.Println(err)
		w.logger.DeleteReceiver(w)
	}
}

func (w *WebsocketLoggerReceiver) ReceiveHistory(history []string) {
	resp := &WebsocketResponse{History: history}

	if err := w.conn.WriteJSON(resp); err != nil {
		log.Println(err)
		w.logger.DeleteReceiver(w)
	}
}
