package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/Shopify/sarama"
)

type Message struct {
    Content string `json:"content"`
}

func main() {
    // Load configuration
    kafkaBroker := os.Getenv("KAFKA_BROKER")
    kafkaTopic := os.Getenv("KAFKA_TOPIC")

    if kafkaBroker == "" {
        kafkaBroker = "localhost:9092"
    }
    if kafkaTopic == "" {
        kafkaTopic = "messages"
    }

    // Create Kafka producer configuration
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true

    // Create Kafka producer
    brokerList := strings.Split(kafkaBroker, ",")
    producer, err := sarama.NewSyncProducer(brokerList, config)
    if err != nil {
        log.Fatalf("Failed to create producer: %s", err)
    }
    defer func() {
        if err := producer.Close(); err != nil {
            log.Printf("Failed to close producer: %s", err)
        }
    }()

    // HTTP handler
    http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
        var msg Message
        if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Produce message to Kafka
        message := &sarama.ProducerMessage{
            Topic: kafkaTopic,
            Value: sarama.StringEncoder(msg.Content),
        }

        partition, offset, err := producer.SendMessage(message)
        if err != nil {
            log.Printf("Failed to send message: %s", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        log.Printf("Message sent to partition %d at offset %d", partition, offset)
        w.WriteHeader(http.StatusAccepted)
    })

    // Start HTTP server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %s", err)
    }
}