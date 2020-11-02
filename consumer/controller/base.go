package controller

import (
	"RabbitMQConsumer/consumer/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Server struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (server *Server) app(w http.ResponseWriter, r *http.Request) {

	var M model.Recieved
	err := json.NewDecoder(r.Body).Decode(&M)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var H []model.Hotel
	var R []model.Room
	var RP []model.Rateplan

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

	server.Router = mux.NewRouter()
	server.InitializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("\nListening to port", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
