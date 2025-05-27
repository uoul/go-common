package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/uoul/go-common/log"
	"github.com/uoul/go-common/messaging"
)

func main() {
	logger := log.NewConsoleLogger(log.TRACE)
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	rabbitMq := messaging.NewRabbitMqMessenger(
		ctx,
		logger,
		"localhost",
		5672,
		"guest",
		"guest",
	)

	go produce(ctx, rabbitMq)
	go consume1(ctx, rabbitMq)
	go consume2(ctx, rabbitMq)
	<-ctx.Done()
}

func consume1(ctx context.Context, messenger messaging.IMessenger[messaging.RabbitMqExchange, amqp091.Delivery]) {
	sub := messenger.Subscribe(messaging.RabbitMqExchange{
		Exchange:   "WASINET",
		RoutingKey: "",
		Type:       "topic",
	})
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sub:
			fmt.Println(string(msg.Result.Body))
		}
	}
}

func consume2(ctx context.Context, messenger messaging.IMessenger[messaging.RabbitMqExchange, amqp091.Delivery]) {
	sub := messenger.Subscribe(messaging.RabbitMqExchange{
		Exchange:   "WASINET2",
		RoutingKey: "",
		Type:       "fanout",
	})
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sub:
			fmt.Println(string(msg.Result.Body))
		}
	}
}

func produce(ctx context.Context, messenger messaging.IMessenger[messaging.RabbitMqExchange, amqp091.Delivery]) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messenger.Publish(messaging.RabbitMqExchange{
				Exchange:   "WASINET",
				RoutingKey: "",
				Type:       "topic",
			}, "SOME MESSAGE")
			time.Sleep(5 * time.Second)
		}
	}
}
