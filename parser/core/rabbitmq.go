package core

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type RabbitMQ struct {
	Conn         *amqp.Connection
	Channel      *amqp.Channel
	ListenQueue  *amqp.Queue
	DeliverQueue *amqp.Queue
}

func NewRabbitMQ(conn *amqp.Connection, ch *amqp.Channel, listen, deliver *amqp.Queue) *RabbitMQ {
	return &RabbitMQ{Conn: conn, Channel: ch, ListenQueue: listen, DeliverQueue: deliver}
}

func InitMQ() *RabbitMQ {
	uri := os.Getenv("RABBITMQ_URI")
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	startQueue, err := ch.QueueDeclare(
		"parser-start",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	resultQueue, err := ch.QueueDeclare(
		"parser-result",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	mq := NewRabbitMQ(conn, ch, &startQueue, &resultQueue)

	return mq
}
