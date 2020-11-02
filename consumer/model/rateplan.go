package model

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

//Rateplan structure
type Rateplan struct {
	Hotell             Hotel                `gorm:"not null" json:"hotel"`
	HotelID            string               `gorm:"not null;" json:"hotel_id"`
	Rateplan           string               `gorm:"not null;" json:"rate_plan_id"`
	CancellationPolicy []CancellationPolicy `gorm:"not null" json:"cancellation_policy"`
	Name               string               `gorm:"not null;" json:"name"`
	Conditions         pq.StringArray       `gorm:"not null" json:"other_conditions"`
	MealPlan           string               `gorm:"not null;" json:"meal_plan"`
}

//CancellationPolicy structure
type CancellationPolicy struct {
	Type            string `gorm:"not null;" json:"type"`
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
