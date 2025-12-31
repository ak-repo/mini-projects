package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func main() {

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "events",
		Balancer: &kafka.Hash{},
	})

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {

		var e Event
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		msg := kafka.Message{
			Key:   []byte(e.Type),
			Value: []byte(e.Data),
			Time:  time.Now(),
		}

		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("event publish"))

	})

	log.Println("HTTP server started")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
