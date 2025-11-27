package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ak-repo/chat-application/config"
	"github.com/ak-repo/chat-application/gen/chatpb"
	"github.com/ak-repo/chat-application/internal/adapter/postgres"
	"github.com/ak-repo/chat-application/internal/app"
	"github.com/ak-repo/chat-application/pkg/db"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	// Using the new structure paths
)

func main() {
	// --- 1. Load Configuration ---

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	if cfg.Database.MaxPoolSize <= 0 {
		cfg.Database.MaxPoolSize = 10
	}
	if cfg.Database.MinPoolSize < 1 {
		cfg.Database.MinPoolSize = 1
	}

	pgDB, err := db.NewPostgresDB(context.Background(), cfg)
	if err != nil {
		log.Fatal("failed to connect db:", zap.Error(err))
	}
	defer pgDB.Close()

	// b. Redis
	// Redis
	rdbAddr := cfg.Redis.Host + ":" + cfg.Redis.Port
	rdb := redis.NewClient(&redis.Options{Addr: rdbAddr})

	// --- 3. Dependency Injection (DI) ---
	repo := postgres.NewChatRepo(pgDB.Pool)
	svc := app.NewServer(repo, rdb)

	// --- 4. Start Background Tasks ---
	// Start the Redis listener in a goroutine
	ctx := context.Background()

	go svc.StartRedisSubscriber(ctx)

	// --- 5. Start gRPC Server ---
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Services.Chat.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Services.Chat.Port, err)
	}

	grpcServer := grpc.NewServer()
	chatpb.RegisterChatServiceServer(grpcServer, svc)

	// 6. Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received shutdown signal. Gracefully stopping gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("✅ Chat Service listening on %s (gRPC)", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
	log.Println("Chat Service shut down successfully.")
}
