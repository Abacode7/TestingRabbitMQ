package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main(){
	log.Println("Starting server...")
	go server()

	client()
}

func server() {
	conn, channel, queue := getQueue()
	defer conn.Close()
	defer channel.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body: []byte("Hello RabbitMQ"),
	}

	for {
		err := channel.Publish("", queue.Name, false, false, msg)
		if err != nil {
			log.Println("Error publishing message")
		}
	}
}

func client() {
	conn, channel, queue := getQueue()
	defer conn.Close()
	defer channel.Close()

	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Println("Error consuming messages")
	}

	for msg := range msgs {
		fmt.Printf("Received messages: %s", msg.Body)
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, amqp.Queue){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error connecting to message queue")

	ch, err := conn.Channel()
	failOnError(err, "Error generating channel")

	queue, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnError(err, "Error generating queue")

	return conn, ch, queue
}

func failOnError(err error, message string){
	if err != nil {
		log.Fatalln(message)
	}
}


