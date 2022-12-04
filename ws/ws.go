package ws

import (
	"github.com/gorilla/websocket"
	"github.com/k0tletka/websocket_logger/logger"
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
		_ = w.conn.Close()
	}
}

func (w *WebsocketLoggerReceiver) ReceiveHistory(history []string) {
	resp := &WebsocketResponse{History: history}

	if err := w.conn.WriteJSON(resp); err != nil {
		log.Println(err)
		w.logger.DeleteReceiver(w)
		_ = w.conn.Close()
	}
}
