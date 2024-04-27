package stream

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func ProduceKafkaMessage(topic, msg string) {
	mechanism, _ := scram.Mechanism(scram.SHA256, os.Getenv("KAFKA_USERNAME"), os.Getenv("KAFKA_PASSWORD"))
	w := kafka.Writer{
		Addr:  kafka.TCP(os.Getenv("KAFKA_URL")),
		Topic: topic,
		Transport: &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		},
	}
	w.WriteMessages(context.Background(), kafka.Message{Value: []byte(msg)})
	w.Close()
}

func ConsumeKafkaMessages(topic string, groupId string) {
	mechanism, _ := scram.Mechanism(scram.SHA512, os.Getenv("KAFKA_USERNAME"), os.Getenv("KAFKA_PASSWORD"))
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_URL")},
		GroupID: groupId,
		Topic:   topic,
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30) // Increase the timeout
	defer cancel()
	for {
		message, _ := r.ReadMessage(ctx)
		fmt.Println(message.Partition, message.Offset, string(message.Value))
	}
	r.Close()
}
