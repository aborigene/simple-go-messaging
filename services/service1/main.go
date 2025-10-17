package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Message struct {
    Content string `json:"content"`
}

func main() {
    // Load configuration
    kafkaBroker := os.Getenv("KAFKA_BROKER")
    kafkaTopic := os.Getenv("KAFKA_TOPIC")

    // Create Kafka producer
    producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
    if err != nil {
        log.Fatalf("Failed to create producer: %s", err)
    }
    defer producer.Close()

    // HTTP handler
    http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
        var msg Message
        if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Produce message to Kafka
        deliveryChan := make(chan kafka.Event)
        err := producer.Produce(&kafka.Message{
            TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
            Value:          []byte(msg.Content),
        }, deliveryChan)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Wait for delivery report
        go func() {
            defer close(deliveryChan)
            e := <-deliveryChan
            m := e.(*kafka.Message)
            if m.TopicPartition.Error != nil {
                log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
            } else {
                log.Printf("Delivered message to %v\n", m.TopicPartition)
            }
        }()

        w.WriteHeader(http.StatusAccepted)
    })

    // Start HTTP server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %s", err)
    }
}