package handler

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type LogMsgWebSocketHandler struct {
	logMsgChannel chan string
	wsConnList    map[*websocket.Conn]bool
}

func (it *LogMsgWebSocketHandler) wsHandler(wsConn *websocket.Conn) {
	log.Printf("A new websocket client connected, %s, %s", wsConn.LocalAddr(), wsConn.RemoteAddr())
	it.wsConnList[wsConn] = true
	buf := make([]byte, 1024)
	for {
		_, err := wsConn.Read(buf)
		if err != nil {
			log.Printf("Error reading from WebSocket connection: %v", err)
			break // Exit the loop if there's an error
		}
	}

	defer delete(it.wsConnList, wsConn)
	defer wsConn.Close()
}

func (it *LogMsgWebSocketHandler) broadcastLogMsg() {
	for {
		logMsg := <-it.logMsgChannel
		for conn := range it.wsConnList {
			conn.Write([]byte(logMsg))
		}
	}
}

func NewLogMsgWebSocketHandler(channel chan string) http.Handler {
	h := &LogMsgWebSocketHandler{logMsgChannel: channel, wsConnList: make(map[*websocket.Conn]bool)}
	go h.broadcastLogMsg()
	return h
}

func (it *LogMsgWebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(it.wsHandler).ServeHTTP(w, r)
}
