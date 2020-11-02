package model

//Recieved structure
type Recieved struct {
	OfferID string   `json:"cm_offer_id"`
	Offers  []Offers `json:"offers"`
}

//Offers structure
type Offers struct {
	Hotel    Hotel    `json:"hotel"`
	Room     Room     `json:"room"`
	Rateplan Rateplan `json:"rate_plan"`
}
