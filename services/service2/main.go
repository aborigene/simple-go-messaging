package main

import (
    "context"
    "encoding/json"
    "log"
    "os"
    "os/signal"
    "syscall"
    "github.com/segmentio/kafka-go"
)

type Message struct {
    Data string `json:"data"`
}

func main() {
    // Load configuration from environment variables
    kafkaBroker := os.Getenv("KAFKA_BROKER")
    kafkaTopic := os.Getenv("KAFKA_TOPIC")

    if kafkaBroker == "" {
        kafkaBroker = "localhost:9092"
    }
    if kafkaTopic == "" {
        kafkaTopic = "messages"
    }

    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{kafkaBroker},
        Topic:   kafkaTopic,
        GroupID: "service2-group",
    })

    defer r.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        sigs := make(chan os.Signal, 1)
        signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
        <-sigs
        cancel()
    }()

    for {
        m, err := r.ReadMessage(ctx)
        if err != nil {
            log.Printf("error while reading message: %v", err)
            continue
        }

        var message Message
        if err := json.Unmarshal(m.Value, &message); err != nil {
            log.Printf("error while unmarshalling message: %v", err)
            continue
        }

        log.Printf("received message: %s", message.Data)
    }
}