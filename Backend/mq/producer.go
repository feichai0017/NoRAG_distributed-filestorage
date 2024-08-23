package mq

import (
	"cloud_distributed_storage/config"
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel

func initChannel() bool {
	if channel != nil {
		return true
	}
	var err error

	conn, err = amqp.Dial(config.RabbitURL)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func Publish(exchange, routingKey string, msg []byte) bool {
	if !initChannel() {
		log.Println("Failed to initialize channel")
		return false
	}
	err := channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
