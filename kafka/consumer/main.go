package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "events",
		GroupID: "event-processors",
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("partition=%d offser=%d key=%s value=%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
}
