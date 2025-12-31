package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ak-repo/go-chat-system/internal/config"
	"github.com/ak-repo/go-chat-system/internal/repository/postgres"
	redis_pkg "github.com/ak-repo/go-chat-system/internal/repository/redis"
	"github.com/ak-repo/go-chat-system/internal/service"
	"github.com/ak-repo/go-chat-system/internal/transport/grpc"
	httpHandler "github.com/ak-repo/go-chat-system/internal/transport/http"
	"github.com/ak-repo/go-chat-system/internal/transport/websocket"
	"github.com/ak-repo/go-chat-system/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize logger
	log := logger.New()
	log.Info("Starting chat application...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", "error", err)
	}

	// Initialize database
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err.Error())
	}
	defer db.Close()

	// Initialize Redis
	redisClient := initRedis(cfg)
	defer redisClient.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	messageRepo := postgres.NewMessageRepository(db)
	conversationRepo := postgres.NewConversationRepository(db)
	presenceRepo := redis_pkg.NewPresenceRepository(redisClient)
	cacheRepo := redis_pkg.NewCacheRepository(redisClient)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	chatService := service.NewChatService(messageRepo, conversationRepo, userRepo)
	presenceService := service.NewPresenceService(presenceRepo)
	signalingService := service.NewSignalingService(cacheRepo)
	signalingService.GetAnswer("101")

	// Initialize WebSocket hub
	hub := websocket.NewHub(chatService, presenceService)
	go hub.Run()

	// Initialize HTTP server
	httpHandler := httpHandler.NewHandler(authService, chatService, hub, cfg)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Initialize gRPC server
	grpcServer := grpc.NewServer(chatService, hub)
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatal("Failed to listen for gRPC", "error", err)
	}

	// Start servers
	go func() {
		log.Info("HTTP server starting", "port", cfg.Server.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error", "error", err)
		}
	}()

	go func() {
		log.Info("gRPC server starting", "port", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal("gRPC server error", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("HTTP server shutdown error", "error", err)
	}

	grpcServer.GracefulStop()
	log.Info("Servers stopped")
}

func initDB(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.Name)

	log.Println(dsn)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func initRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}
