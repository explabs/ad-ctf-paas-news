package main

import (
	"encoding/json"
	"fmt"
	"github.com/explabs/ad-ctf-paas-exploits/service/sender"
	"github.com/explabs/ad-ctf-paas-exploits/service/storage"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Message struct {
	Round string `json:"round"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	var news storage.NewsStruct
	news.Load()
	fmt.Println(news)
	uploadErr := news.UploadNews()
	if uploadErr != nil {
		log.Fatal(uploadErr)
	}

	var host, port = "rabbitmq", 5672
	if os.Getenv("MODE") == "dev" {
		host = "localhost"
	}
	rabbitAddr := fmt.Sprintf("amqp://service:%s@%s:%d", os.Getenv("ADMIN_PASS"), host, port)

	conn, err := amqp.Dial(rabbitAddr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"news", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var m Message
			log.WithFields(log.Fields{"body": string(d.Body)}).Info("Received a message")
			err := json.Unmarshal(d.Body, &m)
			if err != nil {
				log.Error(err)
			}
			round, err := strconv.Atoi(m.Round)
			if err != nil {
				return
			}

			err = sender.SendNews(round)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	go log.Info("News service started")
	<-forever
}
