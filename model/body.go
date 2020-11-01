package model

import "fmt"

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

func (o *Offers) PrintValues() {

	fmt.Println(o)
}
