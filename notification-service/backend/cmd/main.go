package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/segmentio/kafka-go"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- CONFIG ---
var (
	DB_DSN      = "host=localhost user=user password=password dbname=notifications_db port=5433 sslmode=disable"
	KAFKA_BROKER = "localhost:9093" // Connect to host port
	KAFKA_TOPIC  = "domain.events"
	KAFKA_GROUP  = "notification-service-group"
	SERVER_PORT  = ":8080"
	FCM_KEY_PATH = "./serviceAccountKey.json"
)

// --- MODELS ---
type Notification struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    string          `gorm:"index" json:"user_id"`
	Title     string          `json:"title"`
	Body      string          `json:"body"`
	Data      json.RawMessage `gorm:"type:jsonb" json:"data"`
	IsRead    bool            `json:"is_read"`
	SentFCM   bool            `json:"sent_fcm"`
	CreatedAt time.Time       `json:"created_at"`
}

type DeviceToken struct {
	UserID    string    `gorm:"primaryKey" json:"user_id"`
	Token     string    `gorm:"primaryKey" json:"token"`
	Platform  string    `json:"platform"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventPayload struct {
	EventID   string          `json:"event_id"`
	Type      string          `json:"type"`
	UserID    string          `json:"user_id"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}

// --- GLOBALS ---
var (
	db    *gorm.DB
	fcm   *messaging.Client
	wsHub *Hub
)

// --- WEBSOCKET HUB ---
type Hub struct {
	clients map[string][]*websocket.Conn
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string][]*websocket.Conn)}
}

func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = append(h.clients[userID], conn)
	log.Printf("[WS] User %s connected. Active tabs: %d", userID, len(h.clients[userID]))
}

func (h *Hub) Unregister(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.clients[userID]
	for i, c := range conns {
		if c == conn {
			h.clients[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

func (h *Hub) SendJSON(userID string, payload interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conns, ok := h.clients[userID]
	if !ok {
		return
	}
	for _, conn := range conns {
		if err := conn.WriteJSON(payload); err != nil {
			log.Printf("[WS] Write Error: %v", err)
		}
	}
}

// --- KAFKA CONSUMER ---
func startKafkaConsumer(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{KAFKA_BROKER},
		Topic:   KAFKA_TOPIC,
		GroupID: KAFKA_GROUP,
		MinBytes: 10e3,   // 10KB
		MaxBytes: 10e6,   // 10MB
	})
	defer r.Close()

	log.Println("[Kafka] Consumer started...")

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			log.Printf("[Kafka] Fetch Error: %v", err)
			time.Sleep(time.Second)
			continue
		}
		processMessage(ctx, m)
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Printf("[Kafka] Commit Error: %v", err)
		}
	}
}

func processMessage(ctx context.Context, m kafka.Message) {
	var event EventPayload
	if err := json.Unmarshal(m.Value, &event); err != nil {
		log.Printf("[Kafka] Unmarshal Error: %v", err)
		return
	}

	notif := Notification{
		UserID:    event.UserID,
		Title:     "New Notification",
		Body:      fmt.Sprintf("Event type: %s", event.Type),
		Data:      event.Data,
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	if event.Type == "file_uploaded" {
		notif.Title = "File Uploaded"
		notif.Body = "Your file processing is complete."
	}

	if err := db.Create(&notif).Error; err != nil {
		log.Printf("[DB] Save Error: %v", err)
		return
	}
	log.Printf("[Logic] Notification saved for user %s (ID: %d)", notif.UserID, notif.ID)

	wsHub.SendJSON(notif.UserID, notif)
	go sendFCM(notif)
}

func sendFCM(notif Notification) {
	if fcm == nil {
		return
	}
	var tokens []DeviceToken
	db.Where("user_id = ?", notif.UserID).Find(&tokens)
	if len(tokens) == 0 {
		return
	}
	for _, t := range tokens {
		msg := &messaging.Message{
			Token: t.Token,
			Notification: &messaging.Notification{
				Title: notif.Title,
				Body:  notif.Body,
			},
			Data: map[string]string{"notification_id": fmt.Sprintf("%d", notif.ID)},
		}
		if _, err := fcm.Send(context.Background(), msg); err != nil {
			log.Printf("[FCM] Send Error for %s: %v", t.Token[:10], err)
		} else {
			log.Printf("[FCM] Sent successfully to %s", t.Token[:10])
		}
	}
	db.Model(&notif).Update("sent_fcm", true)
}

// --- MAIN ---
func main() {
	var err error
	db, err = gorm.Open(postgres.Open(DB_DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("[Init] DB Connection failed:", err)
	}
	db.AutoMigrate(&Notification{}, &DeviceToken{})
	log.Println("[Init] Database connected.")

	if _, err := os.Stat(FCM_KEY_PATH); err == nil {
		opt := option.WithCredentialsFile(FCM_KEY_PATH)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Printf("[Init] Firebase Error: %v", err)
		} else {
			fcm, _ = app.Messaging(context.Background())
			log.Println("[Init] Firebase FCM initialized.")
		}
	} else {
		log.Println("[Init] No serviceAccountKey.json found. FCM disabled.")
	}

	wsHub = NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	go startKafkaConsumer(ctx)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/api/tokens", func(c *fiber.Ctx) error {
		var input DeviceToken
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		input.UpdatedAt = time.Now()
		if err := db.Save(&input).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(200)
	})

	app.Get("/api/notifications", func(c *fiber.Ctx) error {
		userID := c.Query("user_id")
		var list []Notification
		db.Where("user_id = ?", userID).Order("created_at desc").Limit(50).Find(&list)
		return c.JSON(list)
	})

	app.Patch("/api/notifications/:id/read", func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Model(&Notification{}).Where("id = ?", id).Update("is_read", true)
		return c.SendStatus(200)
	})

	app.Post("/api/debug/publish", func(c *fiber.Ctx) error {
		var payload EventPayload
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		if payload.EventID == "" {
			payload.EventID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		if payload.Timestamp.IsZero() {
			payload.Timestamp = time.Now()
		}

		w := &kafka.Writer{
			Addr:     kafka.TCP(KAFKA_BROKER),
			Topic:    KAFKA_TOPIC,
			Balancer: &kafka.LeastBytes{},
		}
		defer w.Close()

		val, _ := json.Marshal(payload)
		if err := w.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(payload.UserID),
			Value: val,
		}); err != nil {
			return c.Status(500).SendString("Kafka write failed: " + err.Error())
		}
		return c.JSON(fiber.Map{"status": "published", "payload": payload})
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/notifications", websocket.New(func(c *websocket.Conn) {
		userID := c.Query("user_id")
		if userID == "" {
			return
		}
		wsHub.Register(userID, c)
		defer wsHub.Unregister(userID, c)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	go func() {
		if err := app.Listen(SERVER_PORT); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	cancel()
	app.Shutdown()
}
