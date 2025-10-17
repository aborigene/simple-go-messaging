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
    config, err := os.ReadFile("config.yaml")
    if err != nil {
        log.Fatalf("failed to read config file: %v", err)
    }

    var kafkaConfig struct {
        Broker string `json:"broker"`
        Topic  string `json:"topic"`
    }
    if err := json.Unmarshal(config, &kafkaConfig); err != nil {
        log.Fatalf("failed to unmarshal config: %v", err)
    }

    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{kafkaConfig.Broker},
        Topic:   kafkaConfig.Topic,
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