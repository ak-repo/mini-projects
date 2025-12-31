package websocket

import (
	"log"
	"time"

	domain "github.com/ak-repo/go-chat-system/internal/domian"
	"github.com/fasthttp/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512KB
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan *Message
	userID string
}

type Message struct {
	Type       string                 `json:"type"`
	Payload    map[string]interface{} `json:"payload"`
	Recipients []string               `json:"-"`
}

func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan *Message, 256),
		userID: userID,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.hub.presenceServ.SetOnline(c.userID)
		return nil
	})

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		c.handleMessage(&msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case "chat_message":
		c.handleChatMessage(msg)
	case "typing":
		c.handleTyping(msg)
	case "read_receipt":
		c.handleReadReceipt(msg)
	case "call_offer":
		c.handleCallOffer(msg)
	case "call_answer":
		c.handleCallAnswer(msg)
	case "ice_candidate":
		c.handleICECandidate(msg)
	}
}

func (c *Client) handleChatMessage(msg *Message) {
	conversationID, _ := msg.Payload["conversation_id"].(string)
	content, _ := msg.Payload["content"].(string)
	messageType, _ := msg.Payload["message_type"].(string)

	if conversationID == "" || content == "" {
		return
	}

	var msgType domain.MessageType
	switch messageType {
	case "text":
		msgType = domain.MessageTypeText
	case "image":
		msgType = domain.MessageTypeImage
	case "file":
		msgType = domain.MessageTypeFile
	default:
		msgType = domain.MessageTypeText
	}

	message, err := c.hub.chatService.SendMessage(c.userID, conversationID, content, msgType)
	if err != nil {
		c.sendError(err.Error())
		return
	}

	// Get conversation members
	members, err := c.hub.chatService.GetConversationMembers(conversationID)
	if err != nil {
		return
	}

	// Broadcast to online members
	recipients := []string{}
	for _, member := range members {
		if member.UserID != c.userID {
			recipients = append(recipients, member.UserID)
		}
	}

	broadcastMsg := &Message{
		Type: "new_message",
		Payload: map[string]interface{}{
			"id":              message.ID,
			"conversation_id": message.ConversationID,
			"sender_id":       message.SenderID,
			"content":         message.Content,
			"message_type":    message.MessageType,
			"created_at":      message.CreatedAt,
		},
		Recipients: recipients,
	}

	c.hub.broadcast <- broadcastMsg
}

func (c *Client) handleTyping(msg *Message) {
	conversationID, _ := msg.Payload["conversation_id"].(string)
	isTyping, _ := msg.Payload["is_typing"].(bool)

	// Get conversation members
	members, _ := c.hub.chatService.GetConversationMembers(conversationID)
	recipients := []string{}
	for _, member := range members {
		if member.UserID != c.userID {
			recipients = append(recipients, member.UserID)
		}
	}

	typingMsg := &Message{
		Type: "typing_indicator",
		Payload: map[string]interface{}{
			"conversation_id": conversationID,
			"user_id":         c.userID,
			"is_typing":       isTyping,
		},
		Recipients: recipients,
	}

	c.hub.broadcast <- typingMsg
}

func (c *Client) handleReadReceipt(msg *Message) {
	messageID, _ := msg.Payload["message_id"].(string)
	if messageID == "" {
		return
	}

	c.hub.chatService.MarkAsRead(messageID, c.userID)
}

func (c *Client) handleCallOffer(msg *Message) {
	callID, _ := msg.Payload["call_id"].(string)
	recipientID, _ := msg.Payload["recipient_id"].(string)
	offer, _ := msg.Payload["offer"].(string)

	if callID == "" || recipientID == "" || offer == "" {
		return
	}

	offerMsg := &Message{
		Type: "call_offer",
		Payload: map[string]interface{}{
			"call_id":   callID,
			"caller_id": c.userID,
			"offer":     offer,
		},
		Recipients: []string{recipientID},
	}

	c.hub.broadcast <- offerMsg
}

func (c *Client) handleCallAnswer(msg *Message) {
	callID, _ := msg.Payload["call_id"].(string)
	callerID, _ := msg.Payload["caller_id"].(string)
	answer, _ := msg.Payload["answer"].(string)

	answerMsg := &Message{
		Type: "call_answer",
		Payload: map[string]interface{}{
			"call_id": callID,
			"answer":  answer,
		},
		Recipients: []string{callerID},
	}

	c.hub.broadcast <- answerMsg
}

func (c *Client) handleICECandidate(msg *Message) {
	recipientID, _ := msg.Payload["recipient_id"].(string)
	candidate, _ := msg.Payload["candidate"]

	iceMsg := &Message{
		Type: "ice_candidate",
		Payload: map[string]interface{}{
			"sender_id": c.userID,
			"candidate": candidate,
		},
		Recipients: []string{recipientID},
	}

	c.hub.broadcast <- iceMsg
}

func (c *Client) sendError(errMsg string) {
	msg := &Message{
		Type: "error",
		Payload: map[string]interface{}{
			"message": errMsg,
		},
	}
	c.send <- msg
}
