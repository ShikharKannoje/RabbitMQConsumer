package main

import (
	//	"consumer/model"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
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
	var M Recieved
	var H []Hotel
	var R []Room
	var RP []Rateplan
	forever := make(chan bool)
	//dd := make(map[string]interface{})
	go func() {
		for d := range msgs {
			//fmt.Printf("Recieved Message: %s\n", d.Body)
			err = json.Unmarshal(d.Body, &M)
			if err != nil {
				fmt.Println(err)
				panic(err)
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

			fmt.Println("Hotel Details", H)
			fmt.Println("Room Details", R)
			fmt.Println("RatePlan Details", RP)

		}
	}()

	fmt.Println("Successfully connected to Rabbit")
	fmt.Println("[*] - waiting for messages")
	<-forever
}

//********************************************Model********************************************************************

//body.go**************************
//Recieved structure
type Recieved struct {
	Offers []Offers `json:"offers"`
}

//Offers structure
type Offers struct {
	Hotel    Hotel    `json:"hotel"`
	Room     Room     `json:"room"`
	Rateplan Rateplan `json:"rate_plan"`
}

//body.go**************************
//hotel.go************************

//Hotel structure
type Hotel struct {
	Hotelid     string      `gorm:"primary_key" json:"hotel_id"`
	Name        string      `gorm:"size:255;not null;unique" json:"name"`
	Country     string      `gorm:"size:255;not null" json:"country"`
	Address     string      `gorm:"size:255;not null" json:"address"`
	Latitude    json.Number `gorm:"not null" json:"latitude"`
	Longitude   json.Number `gorm:"not null" json:"longitude"`
	Telephone   string      `gorm:"size:15;not null" json:"telephone"`
	Amenities   []string    `gorm:"not null" json:"amenities"`
	Description string      `gorm:"size:255;not null" json:"description"`
	RoomCount   json.Number `gorm:"not null" json:"room_count"`
	Currency    string      `gorm:"size:5;not null" json:"currency"`
}

//Prepare prepares before saving into db
func (h *Hotel) Prepare() {
	h.Hotelid = html.EscapeString(strings.TrimSpace(h.Hotelid))
	h.Name = html.EscapeString(strings.TrimSpace(h.Name))
	h.Country = html.EscapeString(strings.TrimSpace(h.Country))
	h.Address = html.EscapeString(strings.TrimSpace(h.Address))
	h.Telephone = html.EscapeString(strings.TrimSpace(h.Telephone))
	h.Description = html.EscapeString(strings.TrimSpace(h.Description))
	h.Currency = html.EscapeString(strings.TrimSpace(h.Currency))

}

//Validate checks if some value is missing
func (h *Hotel) Validate() error {

	if h.Hotelid == "" {
		return errors.New("Required Hotel ID")
	}
	if h.Name == "" {
		return errors.New("Required Hotel Name")
	}
	if h.Country == "" {
		return errors.New("Required Hotel Country")
	}
	if h.Address == "" {
		return errors.New("Required Address")
	}

	if h.Latitude == "" {
		return errors.New("Required Latitude")
	}

	if h.Longitude == "" {
		return errors.New("Required Longitude")
	}

	if h.RoomCount == "" {
		return errors.New("Required Room Count")
	}
	if h.Telephone == "" {
		return errors.New("Required Phone number")
	}
	if h.Description == "" {
		return errors.New("Required Hotel Description")
	}
	if h.Currency == "" {
		return errors.New("Required Currency")
	}
	if h.Amenities == nil {
		return errors.New("Required Amenities")
	}
	return nil
}

//SaveHotel saves in db
func (h *Hotel) SaveHotel(db *gorm.DB) (*Hotel, error) {

	var err error
	err = db.Debug().Create(&h).Error
	if err != nil {
		return &Hotel{}, err
	}
	return h, nil
}

//hotel.go*****************************

//rateplan.go******************************

//Rateplan structure
type Rateplan struct {
	Hotell             Hotel                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"hotel"`
	HotelID            string                `gorm:"size:255;not null;" json:"hotel_id"`
	Rateplan           string                `gorm:"size:255;not null;" json:"rate_plan_id"`
	CancellationPolicy []CancellationPolicyS `gorm:"not null" json:"cancellation_policy"`
	Name               string                `gorm:"size:255;not null;" json:"name"`
	Conditions         []string              `gorm:"not null" json:"other_conditions"`
	MealPlan           string                `gorm:"size:255;not null;" json:"meal_plan"`
}

//CancellationPolicy structure
type CancellationPolicyS struct {
	Type            string `gorm:"size:255;not null;" json:"type"`
	ExpireDayBefore int    `gorm:"not null;" json:"expires_days_before"`
}

//Prepare before saving
func (r *Rateplan) Prepare() {
	r.HotelID = html.EscapeString(strings.TrimSpace(r.HotelID))
	r.Rateplan = html.EscapeString(strings.TrimSpace(r.Rateplan))
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
	r.MealPlan = html.EscapeString(strings.TrimSpace(r.MealPlan))
}

//Validate the input
func (r *Rateplan) Validate() error {

	if r.HotelID == "" {
		return errors.New("Required Hotel ID")
	}
	if r.Rateplan == "" {
		return errors.New("Required Rate Plan")
	}
	if r.Name == "" {
		return errors.New("Required Rate Plan Name")
	}
	if r.MealPlan == "" {
		return errors.New("Required Meal Plan")
	}
	if r.Conditions == nil {
		return errors.New("Required Terms and conditions")
	}

	return nil
}

//SaveRateplan saves in db
func (r *Rateplan) SaveRateplan(db *gorm.DB) (*Rateplan, error) {

	var err error
	err = db.Debug().Model(&Rateplan{}).Create(&r).Error
	if err != nil {
		return &Rateplan{}, err
	}
	if r.Rateplan != "" {
		err = db.Debug().Model(&Hotel{}).Where("id = ?", r.HotelID).Take(&r.Hotell).Error
		if err != nil {
			return &Rateplan{}, err
		}
	}
	return r, nil
}

//rateplan.go*****************************

//room.go*******************

//Room structure
type Room struct {
	Hotell      Hotel    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"hotel"`
	HotelID     string   `gorm:"size:255;not null;" json:"hotel_id"`
	RoomID      string   `gorm:"primary_key" json:"room_id"`
	Discription string   `gorm:"size:255;not null" json:"description"`
	Name        string   `gorm:"size:255;not null" json:"name"`
	Capacity    Capacity `gorm:"not null" json:"capacity"`
}

//Capacity structure
type Capacity struct {
	MaxAdult      json.Number `gorm:"not null" json:"max_adults"`
	ExtraChildren json.Number `gorm:"not null" json:"extra_children"`
}

//Prepare prepares before saving into db
func (r *Room) Prepare() {
	r.HotelID = html.EscapeString(strings.TrimSpace(r.HotelID))
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
	r.RoomID = html.EscapeString(strings.TrimSpace(r.RoomID))
	r.Discription = html.EscapeString(strings.TrimSpace(r.Discription))
}

//Validate the input
func (r *Room) Validate() error {

	if r.HotelID == "" {
		return errors.New("Required Hotel ID")
	}
	if r.RoomID == "" {
		return errors.New("Required Room ID")
	}
	if r.Name == "" {
		return errors.New("Required Room Name")
	}
	if r.Discription == "" {
		return errors.New("Required Discription")
	}
	if r.Capacity.MaxAdult == "" && r.Capacity.ExtraChildren == "" {
		return errors.New("Required Room Capacity")
	}

	return nil
}

//SaveRoom saves in db
func (r *Room) SaveRoom(db *gorm.DB) (*Room, error) {

	var err error
	err = db.Debug().Model(&Room{}).Create(&r).Error
	if err != nil {
		return &Room{}, err
	}
	if r.RoomID != "" {
		err = db.Debug().Model(&Hotel{}).Where("id = ?", r.HotelID).Take(&r.Hotell).Error
		if err != nil {
			return &Room{}, err
		}
	}
	return r, nil
}

//**********room.go

//*********response.go

//JSON formater
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

//ERROR formater
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}

//response.go
