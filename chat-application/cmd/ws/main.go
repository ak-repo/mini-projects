package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ak-repo/chat-application/config"
	chat "github.com/ak-repo/chat-application/gen/chatpb"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// upgrader configures WebSocket handshake parameters
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allows connections from all origins for local dev/testing
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	cfg, err := config.Load() // Use your existing config loader
	if err != nil {
		log.Fatalf("Config load failed: %v", err)
	}

	// --- 1. Dial Chat Service gRPC ---
	// WARNING: Using insecure credentials. Use transport.NewTLS() for production.
	chatServiceAddr := fmt.Sprintf("localhost:%s", cfg.Services.Chat.Port)
	grpcConn, err := grpc.Dial(chatServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gRPC dial to Chat Service (%s) failed: %v", chatServiceAddr, err)
	}
	defer grpcConn.Close()

	chatClient := chat.NewChatServiceClient(grpcConn)

	// --- 2. Setup WebSocket Endpoint ---
	http.HandleFunc("/ws", wsHandler(chatClient, cfg.JWT.Secret))

	// --- 3. Start HTTP Listener ---
	log.Println("✅ API Gateway WebSocket bridge listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(chatClient chat.ChatServiceClient, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Upgrade to WebSocket
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		defer ws.Close()

		// 2. Authentication (Via Query Param for simplicity)
		token := r.URL.Query().Get("token")
		if token == "" {
			ws.WriteMessage(websocket.TextMessage, []byte(`{"error":"missing token"}`))
			return
		}

		// 3. Open gRPC Stream to Chat Service
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		// The gRPC stream is bidirectional
		stream, err := chatClient.Stream(ctx)
		if err != nil {
			log.Println("Failed to create gRPC stream:", err)
			ws.WriteMessage(websocket.TextMessage, []byte(`{"error":"internal stream failure"}`))
			return
		}

		// 4. Send Auth Payload (first message in the gRPC stream)
		authReq := &chat.StreamRequest{
			Payload: &chat.StreamRequest_Auth{
				Auth: &chat.Auth{Token: token},
			},
		}
		if err := stream.Send(authReq); err != nil {
			log.Println("Failed to send gRPC auth:", err)
			return
		}

		// 5. Start Goroutines for data proxying
		errc := make(chan error, 2)
		go wsToGrpc(ws, stream, errc)
		go grpcToWs(ws, stream, errc)

		// 6. Wait for error (either side closing/failing)
		if err := <-errc; err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Bridge closed due to error:", err)
			}
		}
	}
}

// wsToGrpc reads messages from the WebSocket client and forwards them to the gRPC stream.
func wsToGrpc(ws *websocket.Conn, stream chat.ChatService_StreamClient, errc chan error) {
	defer func() {
		// Close the gRPC send side when the WS connection closes
		if err := stream.CloseSend(); err != nil && err != io.EOF {
			log.Println("Error closing gRPC send:", err)
		}
	}()

	for {
		// Reads JSON from the client (e.g., {"type": "message", "text": "Hi"})
		_, msg, err := ws.ReadMessage()
		if err != nil {
			errc <- err
			return
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(msg, &raw); err != nil {
			log.Println("WS JSON unmarshal error:", err)
			continue
		}

		// Only forward messages of type "message" as ClientMessage
		if typ, ok := raw["type"].(string); ok && typ == "message" {
			cm := &chat.ClientMessage{
				ClientId:  toStr(raw["client_id"]),
				From:      toStr(raw["from"]),
				To:        toStr(raw["to"]),
				ChatId:    toStr(raw["chat_id"]),
				Text:      toStr(raw["text"]),
				MsgType:   toStr(raw["msg_type"]),
				CreatedAt: int64(toFloat(raw["created_at"])),
			}

			req := &chat.StreamRequest{Payload: &chat.StreamRequest_ClientMessage{ClientMessage: cm}}
			if err := stream.Send(req); err != nil {
				errc <- err
				return
			}
		}
	}
}

// grpcToWs reads messages from the gRPC stream and forwards them to the WebSocket client.
func grpcToWs(ws *websocket.Conn, stream chat.ChatService_StreamClient, errc chan error) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			errc <- err
			return
		}

		// Convert gRPC ServerMessage to a standard JSON format for the browser
		if sm := resp.GetServerMessage(); sm != nil {
			b, _ := json.Marshal(map[string]interface{}{
				"type":       "message",
				"server_id":  sm.ServerId,
				"client_id":  sm.ClientId,
				"from":       sm.From,
				"to":         sm.To,
				"chat_id":    sm.ChatId,
				"text":       sm.Text,
				"msg_type":   sm.MsgType,
				"created_at": sm.CreatedAt,
			})

			if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
				errc <- err
				return
			}
		}
	}
}

// Helper functions for type safety in JSON handling
func toStr(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0.0
}
