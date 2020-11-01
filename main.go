package main

import (
	"consumer/model"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func main() {

	fmt.Println("Consumer app")
	conn, err := amqp.Dial(os.Getenv("RABITTMQ_CRED"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer ch.Close()

	msgs, err := ch.Consume(os.Getenv("QUEUE_NAME"), "", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var M model.Recieved
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//fmt.Printf("Recieved Message: %s\n", d.Body)
			js, _ := json.Marshal(d.Body)
			fmt.Println("JSON Format Here", js)
			err = json.Unmarshal(js, M)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println(M)
		}
	}()

	fmt.Println("Successfully connected to Rabbit")
	fmt.Println("[*] - waiting for messages")
	<-forever
}
