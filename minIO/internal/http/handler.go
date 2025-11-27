package httpapi

import (
    "context"
    "fmt"
    "io"
    "log"
    "minio-demo/internal/storage"
    "net/http"
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/minio/minio-go/v7"
)

type API struct {
    Store *storage.MinioClient
}

func NewAPI(store *storage.MinioClient) *API {
    return &API{Store: store}
}

func (a *API) UploadHandler(c *fiber.Ctx) error {
    fileHeader, err := c.FormFile("file")
    if err != nil {
        return c.Status(http.StatusBadRequest).SendString("missing file: " + err.Error())
    }

    file, err := fileHeader.Open()
    if err != nil {
        return c.Status(http.StatusInternalServerError).SendString(err.Error())
    }
    defer file.Close()

    objectName := c.FormValue("path")
    if objectName == "" {
        objectName = "uploads/" + fileHeader.Filename
    }

    size := fileHeader.Size

    contentType := fileHeader.Header.Get("Content-Type")
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    info, err := a.Store.Client.PutObject(
        context.Background(),
        a.Store.Bucket,
        objectName,
        file,
        size,
        minio.PutObjectOptions{ContentType: contentType},
    )
    if err != nil {
        log.Printf("PutObject error: %v", err)
        return c.Status(500).SendString("upload failed")
    }

    return c.JSON(fiber.Map{
        "ok":          true,
        "object_name": objectName,
        "size":        info.Size,
    })
}

func (a *API) SignedURLHandler(c *fiber.Ctx) error {
    object := c.Query("object")
    if object == "" {
        return c.Status(400).SendString("object required")
    }

    expiry := storage.SignedURLExpiry()
    if q := c.Query("expiry"); q != "" {
        if v, err := strconv.Atoi(q); err == nil {
            expiry = int64(v)
        }
    }

    url, err := a.Store.PresignedGetURL(context.Background(), object, expiry)
    if err != nil {
        log.Printf("presigned url err: %v", err)
        return c.Status(500).SendString("couldn't generate url")
    }

    return c.JSON(fiber.Map{
        "url": url,
    })
}

func (a *API) Ping(c *fiber.Ctx) error {
    return c.SendString("ok")
}

func (a *API) ProxyDownloadHandler(c *fiber.Ctx) error {
    object := c.Query("object")
    if object == "" {
        return c.Status(400).SendString("object required")
    }

    reader, err := a.Store.Client.GetObject(
        context.Background(),
        a.Store.Bucket,
        object,
        minio.GetObjectOptions{}, // FIXED
    )
    if err != nil {
        return c.Status(500).SendString("failed to get object")
    }

    c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", object))
    _, err = io.Copy(c.Context().Response.BodyWriter(), reader)
    if err != nil {
        return c.Status(500).SendString("stream failed: " + err.Error())
    }

    return nil
}
