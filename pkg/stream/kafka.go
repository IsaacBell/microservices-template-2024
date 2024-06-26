package stream

import (
	"context"
	"core/internal/util"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var (
	brokers []string
	sigs    chan os.Signal

	readersMu sync.Mutex
	readers   map[string]*kafka.Reader
	writersMu sync.Mutex
	writers   map[string]*kafka.Writer

	batchSize uint16 = 100

	Once sync.Once
)

type callbackFn func(msg string)

func initKafka() {
	Once.Do(func() {
		brokers = []string{os.Getenv("KAFKA_URL")}
		for _, broker := range brokers {
			if broker == "" {
				log.Fatalln("Kafka Broker Not Found")
			}
		}

		sigs = make(chan os.Signal, 1)
		readers = make(map[string]*kafka.Reader) // 1 reader for each topic

		util.PrintLnInColor(util.AnsiBrightGreen, "Kafka initialized")
		util.PrintLnInColor(util.AnsiBackgroundMagenta, "  Brokers:")
		for _, broker := range brokers {
			util.PrintLnInColor(util.AnsiBackgroundMagenta, "    -> %s, ", broker)
		}
	})
}

func createReader(topic, groupID string) *kafka.Reader {
	initKafka()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
		Dialer: &kafka.Dialer{
			SASLMechanism: saslMechanism(),
			TLS:           &tls.Config{},
		},
		MaxBytes: 10e6, // 10MB
	})
	readers[topic] = reader
	// reader.SetOffset(42)
	return reader
}

func createWriter(topic string) *kafka.Writer {
	initKafka()
	writer := &kafka.Writer{
		Addr:  kafka.TCP(os.Getenv("KAFKA_URL")),
		Topic: topic,
		Transport: &kafka.Transport{
			SASL: saslMechanism(),
			TLS:  &tls.Config{},
		},
	}
	writers[topic] = writer
	return writer
}

func saslMechanism() sasl.Mechanism {
	mechanism, _ := scram.Mechanism(scram.SHA256, os.Getenv("KAFKA_USERNAME"), os.Getenv("KAFKA_PASSWORD"))
	return mechanism
}

func consumeKafkaMessages(ctx context.Context, topic string, groupID string, callback callbackFn) error {
	readersMu.Lock()
	reader, ok := readers[topic]
	if !ok {
		reader = createReader(topic, groupID)
	}

	fmt.Println("Consuming Kafka messages from topic: ", topic)

	readersMu.Unlock()
	defer func() {
		readersMu.Lock()
		delete(readers, topic)
		readersMu.Unlock()
		reader.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3000) // Increase the timeout
	defer cancel()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigs)

	go func() {
		sig := <-sigs
		fmt.Printf("Received signal: %s, canceling context...\n", sig)
		cancel()
	}()

	for {
		fmt.Printf("reader.Stats(): %v\n", reader.Stats())
		message, err := reader.ReadMessage(ctx)
		fmt.Println(message.Partition, message.Offset, string(message.Value))

		if err != nil {
			if err == context.Canceled {
				fmt.Print(".")
				fmt.Println("Reader context cancelled")
				continue
			}
			return err
		}

		log.Println("Kafka Message: ", message)
		fmt.Println(message.Partition, message.Offset, string(message.Value))
		callback(string(message.Value))
	}
}

func ProduceKafkaMessage(topic, msg string) error {
	writersMu.Lock()
	writer, ok := writers[topic]
	if !ok {
		writer = createWriter(topic)
	}
	writersMu.Unlock()

	return writer.WriteMessages(context.Background(), kafka.Message{Value: []byte(msg)})
}

/*
		Full Kafka Example:
		ctx, cancel := StartKafkaConsumer("topic1", "group1", func(msg string) {
	    // Process the message
		})
		if someCondition == true {
			cancel()
		}
*/
func StartKafkaConsumer(topic string, groupID string, callback callbackFn) (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := consumeKafkaMessages(ctx, topic, groupID, callback)
		if err != nil {
			fmt.Printf("Error consuming messages from topic '%s': %v\n", topic, err)
		}
	}()

	return ctx, cancel
}

/*
// Stop a consumer
StopKafkaConsumer(ctx1, "topic1")
StopKafkaConsumer(ctx2, "topic2")
*/
func StopKafkaConsumer(ctx context.Context, topic string) {
	readersMu.Lock()
	reader, ok := readers[topic]
	if ok {
		reader.Close()
		delete(readers, topic)
	}
	readersMu.Unlock()

	ctx.Done()
}

func StopKafka(ctx context.Context) {
	readersMu.Lock()
	for _, reader := range readers {
		reader.Close()
	}
	readers = make(map[string]*kafka.Reader)
	readersMu.Unlock()

	writersMu.Lock()
	for _, writer := range writers {
		writer.Close()
	}
	writers = make(map[string]*kafka.Writer)
	writersMu.Unlock()
}
