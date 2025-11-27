package grpcsvc

import (
    "context"
    "io"
    "io/ioutil"
    "log"
    pb "minio-demo/internal/proto"
    "minio-demo/internal/storage"
    "os"

    "github.com/minio/minio-go/v7"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type UploadServer struct {
    pb.UnimplementedUploadServiceServer
    Store *storage.MinioClient
}

func NewUploadServer(store *storage.MinioClient) *UploadServer {
    return &UploadServer{Store: store}
}

func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
    ctx := context.Background()

    // Temp file on disk (good for large chunks)
    tmpFile, err := ioutil.TempFile("", "upload-*")
    if err != nil {
        return status.Errorf(codes.Internal, "temp file error: %v", err)
    }
    defer func() {
        tmpFile.Close()
        os.Remove(tmpFile.Name())
    }()

    var objectName string
    var totalBytes int64

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return status.Errorf(codes.Internal, "recv error: %v", err)
        }

        switch payload := req.Payload.(type) {

        case *pb.UploadRequest_Info:
            if payload.Info == nil {
                return status.Errorf(codes.InvalidArgument, "file info expected")
            }
            objectName = payload.Info.ObjectName
            if objectName == "" {
                return status.Errorf(codes.InvalidArgument, "object_name required")
            }

        case *pb.UploadRequest_Chunk:
            if payload.Chunk == nil {
                continue
            }
            n, err := tmpFile.Write(payload.Chunk.Data)
            if err != nil {
                return status.Errorf(codes.Internal, "write tmp error: %v", err)
            }
            totalBytes += int64(n)

        default:
            return status.Errorf(codes.InvalidArgument, "unknown payload")
        }
    }

    // Reset file pointer for upload
    tmpFile.Sync()
    tmpFile.Seek(0, io.SeekStart)

    // Upload to MinIO (FIXED)
    _, err = s.Store.Client.PutObject(
        ctx,
        s.Store.Bucket,
        objectName,
        tmpFile,
        totalBytes,
        minio.PutObjectOptions{
            ContentType: "application/octet-stream",
        },
    )
    if err != nil {
        log.Printf("minio putobject error: %v", err)
        return status.Errorf(codes.Internal, "upload failed: %v", err)
    }

    return stream.SendAndClose(&pb.UploadResponse{
        Ok:         true,
        Message:    "uploaded",
        ObjectName: objectName,
    })
}

func (s *UploadServer) GetSignedUrl(ctx context.Context, req *pb.GetSignedUrlRequest) (*pb.GetSignedUrlResponse, error) {
    expiry := int64(req.ExpirySeconds)
    if expiry == 0 {
        expiry = storage.SignedURLExpiry()
    }

    url, err := s.Store.PresignedGetURL(ctx, req.ObjectName, expiry)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "presign error: %v", err)
    }

    return &pb.GetSignedUrlResponse{Url: url}, nil
}
