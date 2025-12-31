package http

import (
	"net/http"

	websocket_pkg "github.com/ak-repo/go-chat-system/internal/transport/websocket"
	"github.com/fasthttp/websocket"
	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configure properly in production
	},
}

func (h *Handler) websocketHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := websocket_pkg.NewClient(h.hub, conn, userID)

	go client.WritePump()
	go client.ReadPump()
}
