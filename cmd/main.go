package main

import (
	"fmt"
	"github.com/brcodingdev/stock-service/internal/pkg/app"
	"github.com/brcodingdev/stock-service/internal/pkg/broker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("stock service starting ...")

	err := godotenv.Load()
	if err != nil {
		log.Panicf("could not load .env file: %s", err)
	}

	rmqHost := os.Getenv("RABBIT_HOST")
	rmqUserName := os.Getenv("RABBIT_USERNAME")
	rmqPassword := os.Getenv("RABBIT_PASSWORD")
	rmqPort := os.Getenv("RABBIT_PORT")

	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		rmqUserName,
		rmqPassword,
		rmqHost,
		rmqPort,
	)

	receiverQueue := os.Getenv("RECEIVER_QUEUE")
	publisherQueue := os.Getenv("PUBLISHER_QUEUE")
	stockServiceUrl := os.Getenv("STOCK_SERVICE_URL")

	if receiverQueue == "" ||
		publisherQueue == "" ||
		stockServiceUrl == "" {
		log.Panicln("required RECEIVER_QUEUE, PUBLISHER_QUEUE, STOCK_SERVICE_URL env vars set")
	}

	handlerStockApp := app.NewStockApp(stockServiceUrl, &http.Client{})
	rabbit, err := broker.NewRabbitMQ(
		dsn,
		receiverQueue,
		publisherQueue,
		handlerStockApp)

	if err != nil {
		log.Panicf("could not initialize RabbitMQ, start chat-service first: %s", err)
	}

	defer rabbit.Close()

	// handle graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		fmt.Println("received interrupt signal. shutting down gracefully...")
		// Perform any cleanup or shutdown logic here

		// close the RabbitMQ connection
		if err := rabbit.Close(); err != nil {
			fmt.Println("could not close RabbitMQ:", err)
		}

		fmt.Println("graceful shutdown completed")
		os.Exit(0)
	}()

	go rabbit.Consume()

	select {}
}
