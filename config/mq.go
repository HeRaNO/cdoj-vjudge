package config

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var MQCh *amqp.Channel
var MQ <-chan amqp.Delivery
var QueueName, ReplyQueueName string

func InitMQ() {
	config := conf.MQ
	QueueName = config.QueueName
	amqpUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.UserName, config.Password,
		config.Host, config.Port)
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		log.Println("[FATAL] Init message queue failed: dial failed")
		panic(err)
	}

	MQCh, err = conn.Channel()
	if err != nil {
		log.Println("[FATAL] Init message queue failed: init channel failed")
		panic(err)
	}
	q, err := MQCh.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		log.Println("[FATAL] Init message queue failed: declare queue failed")
		panic(err)
	}
	err = MQCh.Qos(1, 0, false)
	if err != nil {
		log.Println("[FATAL] Init message queue failed: set Qos failed")
		panic(err)
	}
	ReplyQueueName = q.Name
	MQ, err = MQCh.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Println("[FATAL] Init message queue failed: failed to register a consumer")
		panic(err)
	}
	log.Println("[INFO] Init message queue successfully")
}
