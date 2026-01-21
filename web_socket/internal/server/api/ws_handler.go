package api

import (
	"net/http"

	"ws_test/internal/logger"
	"ws_test/internal/server/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (server *Server) wsHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Log.Print(2, "[ws] ws upgrade error: %v", err)
		return
	}

	logger.Log.Print(2, "request client..")

	client := &ws.Client{
		Hub:  server.hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	server.hub.Register <- client

	go client.WsRead()
	go client.WsWrite()
}
