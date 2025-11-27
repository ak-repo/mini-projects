package main

import (
	"context"
	"fmt"
	"io"
	"log"
	pb "minio-demo/internal/proto"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("usage: go run main.go <file-path> <object-name>")
	}
	filePath := os.Args[1]
	objectName := os.Args[2]

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewUploadServiceClient(conn)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer f.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("upload: %v", err)
	}

	// send FileInfo first
	info := &pb.UploadRequest{
		Payload: &pb.UploadRequest_Info{
			Info: &pb.FileInfo{
				Bucket:      "", // server default bucket used
				ObjectName:  objectName,
				ContentType: "application/octet-stream",
			},
		},
	}
	if err := stream.Send(info); err != nil {
		log.Fatalf("send info: %v", err)
	}

	buf := make([]byte, 32*1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("read: %v", err)
		}
		chunk := &pb.UploadRequest{
			Payload: &pb.UploadRequest_Chunk{
				Chunk: &pb.FileChunk{Data: buf[:n]},
			},
		}
		if err := stream.Send(chunk); err != nil {
			log.Fatalf("send chunk: %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("close and recv: %v", err)
	}
	fmt.Printf("Upload response: ok=%v message=%s object=%s\n", resp.Ok, resp.Message, resp.ObjectName)

	// get presigned url
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	signed, err := client.GetSignedUrl(ctx, &pb.GetSignedUrlRequest{
		ObjectName:    objectName,
		ExpirySeconds: 3600,
	})
	if err != nil {
		log.Fatalf("get signed url: %v", err)
	}
	fmt.Printf("Signed URL: %s\n", signed.Url)
}
