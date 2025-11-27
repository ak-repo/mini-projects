package main

import (
	"fmt"
	"log"
	grpcsvc "minio-demo/internal/grpc"
	httpapi "minio-demo/internal/http"
	pb "minio-demo/internal/proto"
	"minio-demo/internal/storage"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// load .env
	_ = godotenv.Load()

	// create minio client
	store := storage.NewMinioClient()

	// setup fiber
	app := fiber.New()

	api := httpapi.NewAPI(store)

	app.Post("/upload", api.UploadHandler)
	app.Get("/signed-url", api.SignedURLHandler)
	app.Get("/proxy-download", api.ProxyDownloadHandler)
	app.Get("/ping", api.Ping)

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = "0.0.0.0:8080"
	}

	// start gRPC server
	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = "0.0.0.0:50051"
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	uploadSrv := grpcsvc.NewUploadServer(store)
	pb.RegisterUploadServiceServer(grpcServer, uploadSrv)

	// run gRPC in goroutine
	go func() {
		fmt.Printf("gRPC serving at %s\n", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// run http server in goroutine
	go func() {
		fmt.Printf("HTTP serving at %s\n", httpAddr)
		if err := app.Listen(httpAddr); err != nil {
			log.Fatalf("fiber server error: %v", err)
		}
	}()

	// wait for signal to gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	fmt.Println("shutting down...")

	grpcServer.GracefulStop()
	app.Shutdown()
	fmt.Println("stopped")
}
