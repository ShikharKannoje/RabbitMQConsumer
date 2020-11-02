package controller

import (
	"RabbitMQConsumer/consumer/model"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
)

type Server struct {
	DB *gorm.DB
}

func (server *Server) consumer() {

	fmt.Println("\nConsumer Started")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
	var M model.Recieved

	var H []model.Hotel
	var R []model.Room
	var RP []model.Rateplan
	msgs, err := ch.Consume(os.Getenv("QUEUE_NAME"), "", true, false, false, false, nil)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Recieved Data: %s\n", d.Body)
			val, err := json.Marshal(d.Body)
			if err != nil {
				fmt.Println(err)
			}
			err = json.Unmarshal(val, &M)
			if err != nil {
				fmt.Println(err)
			}
			for _, j := range M.Offers {
				err = j.Hotel.Validate()
				if err != nil {
					fmt.Println(err)
					panic(err)
				} else {
					H = append(H, j.Hotel)
				}
				err = j.Room.Validate()
				if err != nil {
					fmt.Println(err)
					panic(err)
				} else {
					R = append(R, j.Room)
				}
				err = j.Rateplan.Validate()
				if err != nil {
					fmt.Println(err)
					panic(err)
				} else {
					RP = append(RP, j.Rateplan)
				}

			}

			for _, j := range H {
				j.Prepare()
				j.SaveHotel(server.DB)
			}
			for _, j := range R {
				j.Prepare()
				j.SaveRoom(server.DB)
			}
			for _, j := range RP {
				j.Prepare()
				j.SaveRateplan(server.DB)
			}
			// fmt.Println("Hotel Details", H)
			// fmt.Println("Room Details", R)
			// fmt.Println("RatePlan Details", RP)

		}
	}()
	fmt.Println("Successfully connected to Rabbit MQ")
	fmt.Println("[*] - Waiting for Messages")
	<-forever

}

func (server *Server) Initialize(Dbdriver, DbHost, DbPort, DbUser, DbName, DbPassword string) {

	var err error
	fmt.Println(DbHost, DbPort, DbUser, DbName, DbPassword)
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	fmt.Println(Dbdriver, DBURL)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	//server.Router = mux.NewRouter()
	server.InitializeConsumer()
}
