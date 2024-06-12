# Kafka Streaming Package

The `stream` package provides a convenient way to interact with Apache Kafka for producing and consuming messages. It offers a simple and intuitive API for managing Kafka readers and writers.

## Features

- Automatic initialization of Kafka brokers
- Creation and management of Kafka readers and writers
- Consuming messages from Kafka topics with a callback function
- Producing messages to Kafka topics
- Graceful shutdown and cleanup of Kafka readers and writers

## Usage

### Initialization

The package automatically initializes the Kafka brokers based on the `KAFKA_URL` environment variable. It also sets up the necessary SASL authentication mechanism using the `KAFKA_USERNAME` and `KAFKA_PASSWORD` environment variables.

### Consuming Messages

To start consuming messages from a Kafka topic, use the `StartKafkaConsumer` function:

```go
ctx, cancel := stream.StartKafkaConsumer("topic1", "group1", func(msg string) {
    // Process the message
})
```

The `StartKafkaConsumer` function takes the topic name, consumer group ID, and a callback function as parameters. It returns a context and a cancel function that can be used to stop the consumer.

To stop a specific consumer, use the `StopKafkaConsumer` function:

```go
stream.StopKafkaConsumer(ctx, "topic1")
```

### Producing Messages

To produce a message to a Kafka topic, use the `ProduceKafkaMessage` function:

```go
err := stream.ProduceKafkaMessage("topic1", "Hello, Kafka!")
if err != nil {
    // Handle the error
}
```

The `ProduceKafkaMessage` function takes the topic name and the message as parameters. It returns an error if the message production fails.

### Graceful Shutdown

To perform a graceful shutdown of all Kafka readers and writers, use the `StopKafka` function:

```go
stream.StopKafka(ctx)
```

The `StopKafka` function takes a context as a parameter and closes all the active readers and writers.

## Configuration

The package requires the following environment variables to be set:

- `KAFKA_URL`: The URL of the Kafka broker(s).
- `KAFKA_USERNAME`: The username for SASL authentication.
- `KAFKA_PASSWORD`: The password for SASL authentication.

Make sure to set these environment variables before running your application.

## Examples

Produce a Kafka message:

```go
import "core/pkg/stream"

func HelloWorld() {
  err := stream.ProduceKafkaMessage("main", "Hello, world!")
}
```

Consume a Kafka topic:

```go
topic := "example"

stream.StartKafkaConsumer(topic, serviceName, func(msg string) {
  fmt.PrintLn("Kafka: [", topic, "] ", msg)
})
```

Use Kafka Consumer in streaming requests:

```go
func (s *LogService) ConsumeStream(req *v1.ConsumeRequest, conn v1.Log_ConsumeStreamServer) error {
	ctx, cancel := stream.StartKafkaConsumer(req.Topic, "core", func(msg string) {
		fmt.Println("Received Kafka Msg: [", req.Topic, "] ", msg)
		err := conn.Send(&v1.ConsumeResponse{
			Record: &v1.Record{
				Value:  []byte(msg),
				Offset: 0,
			},
		})
		if err != nil {
			return
		}
	})
	defer cancel()

	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				cancel()
				return
			default:
				// Do nothing
			}
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}
```

## License

This package is open-source and available under the [MIT License](LICENSE).