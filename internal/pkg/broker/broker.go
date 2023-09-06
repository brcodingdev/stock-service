package broker

import "github.com/brcodingdev/stock-service/internal/pkg/broker/event"

// Broker contract to handle pub/sub messages
// In this case we are using RabbitMQ, but we can use other like Kafka
type Broker interface {
	// Publish publishes events
	Publish(request event.StockResponse)
	// Consume consumes events
	Consume()
	// Close to close connections
	Close() error
}
