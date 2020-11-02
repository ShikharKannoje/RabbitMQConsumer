package model

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

//Room structure
type Room struct {
	Hotell      Hotel    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"hotel"`
	HotelID     string   `gorm:"not null;" json:"hotel_id"`
	RoomID      string   `gorm:"primary_key" json:"room_id"`
	Description string   `gorm:"not null" json:"description"`
	Name        string   `gorm:"not null" json:"name"`
	Capacity    Capacity `gorm:"not null" json:"capacity"`
}

//Capacity structure
type Capacity struct {
	MaxAdult      int `gorm:"not null" json:"max_adults"`
	ExtraChildren int `gorm:"not null" json:"extra_children"`
}

//Prepare prepares before saving into db
func (r *Room) Prepare() {
	r.HotelID = html.EscapeString(strings.TrimSpace(r.HotelID))
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
	r.RoomID = html.EscapeString(strings.TrimSpace(r.RoomID))
	r.Description = html.EscapeString(strings.TrimSpace(r.Description))
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
	if r.Description == "" {
		return errors.New("Required Description")
	}
	if r.Capacity.MaxAdult == 0 && r.Capacity.ExtraChildren == 0 {
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

// "room": {
// 	"hotel_id": "BH~46456",
// 	"room_id": "S2Q",
// 	"description": "JUNIOR SUITES WITH 2 QUEEN BEDS",
// 	"name": "JUNIOR SUITES WITH 2 QUEEN BEDS",
// 	"capacity": {
// 		"max_adults": 2,
// 		"extra_children": 2
// 	}
// },
