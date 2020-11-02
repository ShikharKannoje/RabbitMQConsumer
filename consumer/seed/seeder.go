package seed

import (
	"RabbitMQConsumer/consumer/model"
	"log"

	"github.com/jinzhu/gorm"
)

var hotels = []model.Hotel{
	model.Hotel{
		Hotelid:     "Ritz-22334",
		Name:        "RitzCarlton",
		Country:     "India",
		Address:     "Bangalore",
		Latitude:    -323.4343,
		Longitude:   324.3553,
		Telephone:   "+91-334344343",
		Amenities:   []string{"play area", "pool", "buffet"},
		Description: "Cheap & Best",
		RoomCount:   234,
		Currency:    "INR",
	},
}

var rooms = []model.Room{
	model.Room{
		RoomID:      "Ref33",
		Description: "Hello world 1",
		Name:        "Delux suite",
		Capacity: model.Capacity{
			MaxAdult:      2,
			ExtraChildren: 3,
		},
	},
}

var can = []model.CancellationPolicy{
	model.CancellationPolicy{
		Type:            "hjk",
		ExpireDayBefore: 3,
	},
}

var rateplans = []model.Rateplan{
	model.Rateplan{
		Rateplan:           "Title 1",
		CancellationPolicy: can,
		Name:               "Best",
		Conditions:         []string{"abc", "cde"},
		MealPlan:           "Breakfast only",
	},
}

func Load(db *gorm.DB) {

	for i, _ := range hotels {
		err := db.Debug().Model(&model.Hotel{}).Create(&hotels[i]).Error
		if err != nil {
			log.Fatalf("cannot seed hotel table: %v", err)
		}
		rooms[i].HotelID = hotels[i].Hotelid

		err = db.Debug().Model(&model.Room{}).Create(&rooms[i]).Error
		if err != nil {
			log.Fatalf("cannot seed room table: %v", err)
		}
		rateplans[i].HotelID = hotels[i].Hotelid
		err = db.Debug().Model(&model.Rateplan{}).Create(&rateplans[i]).Error
		if err != nil {
			log.Fatalf("cannot seed rateplan table: %v", err)
		}
	}
}
