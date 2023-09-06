package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brcodingdev/stock-service/internal/pkg/broker/event"
	"log"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ implements Broker RabbitMQ
type RabbitMQ struct {
	receiverQueue  *amqp.Queue
	publisherQueue *amqp.Queue
	channel        *amqp.Channel
	connection     *amqp.Connection
	handler        Handler
}

// NewRabbitMQ ...
func NewRabbitMQ(
	dsn string,
	receiverQueue string,
	publisherQueue string,
	handler Handler,
) (*RabbitMQ, error) {
	conn, err := amqp.Dial(dsn)

	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	receiver, err := channel.QueueDeclare(
		receiverQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	publisher, err := channel.QueueDeclare(
		publisherQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		receiverQueue:  &receiver,
		publisherQueue: &publisher,
		channel:        channel,
		connection:     conn,
		handler:        handler,
	}, nil
}

// Close ...
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return err
		}
	}

	if r.connection != nil {
		if err := r.connection.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Publish ...
func (r *RabbitMQ) Publish(request event.StockResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(request)
	if err != nil {
		fmt.Println("could not marshal request ", err)
		return
	}

	err = r.channel.PublishWithContext(ctx,
		"",
		r.publisherQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	if err != nil {
		fmt.Println("could not publish message ", err)
		return
	}

	log.Printf(" message sent %s\n", body)
}

// Consume ...
func (r *RabbitMQ) Consume() {
	messages, err := r.channel.Consume(
		r.receiverQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("could not consume message %s", err)
		return
	}

	chanReceivedMessages := make(chan event.StockRequest)
	go messageTransformer(messages, chanReceivedMessages)
	go processRequest(chanReceivedMessages, r)
	log.Printf("waiting messages...")
}

func messageTransformer(
	chanReceivedMessages <-chan amqp.Delivery,
	receivedMessages chan event.StockRequest,
) {
	var sr event.StockRequest
	for d := range chanReceivedMessages {
		log.Println("received message", string(d.Body))
		err := json.Unmarshal(d.Body, &sr)
		if err != nil {
			log.Printf("could not receive message %s ", err)
			continue
		}

		log.Println("received a request")
		receivedMessages <- sr
	}
}

func processRequest(request <-chan event.StockRequest, r *RabbitMQ) {
	for req := range request {
		log.Println("stock request", req.ChatRoomID)
		cM := req.ChatMessage
		cM = strings.Replace(cM, "/stock=", "", 1)
		// notify user that request is processing
		processingMsg := event.StockResponse{
			RoomID:  req.ChatRoomID,
			Message: fmt.Sprintf("processing: %s", cM),
		}
		// notify in parallel
		go r.Publish(processingMsg)
		// handle request and notify user
		msg := r.handler.HandleStockRequest(cM)
		responseMsg := event.StockResponse{
			RoomID:  req.ChatRoomID,
			Message: msg,
		}
		// notify in parallel
		go r.Publish(responseMsg)
		log.Println("processed", req.ChatMessage)
	}
}
