package model

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

//Hotel structure
type Hotel struct {
	Hotelid     string   `gorm:"primary_key" json:"hotel_id"`
	Name        string   `gorm:"size:255;not null;unique" json:"name"`
	Country     string   `gorm:"size:255;not null" json:"country"`
	Address     string   `gorm:"size:255;not null" json:"address"`
	Latitude    float64  `gorm:"not null" json:"latitude"`
	Longitude   float64  `gorm:"not null" json:"longitude"`
	Telephone   string   `gorm:"size:15;not null" json:"telephone"`
	Amenities   []string `gorm:"not null" json:"amenities"`
	Description string   `gorm:"size:255;not null" json:"discription"`
	RoomCount   int      `gorm:"not null" json:"room_count"`
	Currency    string   `gorm:"size:5;not null" json:"currency"`
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

	if h.Latitude == 0 {
		return errors.New("Required Latitude")
	}

	if h.Longitude == 0 {
		return errors.New("Required Longitude")
	}

	if h.RoomCount == 0 {
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

// "hotel": {
// 	"hotel_id": "BH~46456",
// 	"name": "Hawthorn Suites by Wyndham Eagle CO",
// 	"country": "US",
// 	"address": "0315 Chambers Avenue, 81631",
// 	"latitude": 39.660193,
// 	"longitude": -106.824123,
// 	"telephone": "+1-970-3283000",
// 	"amenities": [
// 		"Business Centre",
// 		"Fitness Room/Gym",
// 		"Pet Friendly",
// 		"Disabled Access",
// 		"Air Conditioned",
// 		"Free WIFI",
// 		"Elevator / Lift",
// 		"Parking"
// 	],
// 	"description": "Stay a while in beautiful mountain country at this Hawthorn Suites by Wyndham Eagle CO hotel, just off Interstate 70, only 6 miles from the Vail/Eagle Airport and close to skiing, golfing, Eagle River and great restaurants. Pets are welcome at this h",
// 	"room_count": 1,
// 	"currency": "USD"
// }
