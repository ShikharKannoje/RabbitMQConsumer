package controller

import (
	"RabbitMQConsumer/consumer/model"
	"RabbitMQConsumer/consumer/response"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (server *Server) CreateHotel(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	hotel := model.Hotel{}
	err = json.Unmarshal(body, &hotel)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	hotel.Prepare()
	err = hotel.Validate()
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	hotelCreated, err := hotel.SaveHotel(server.DB)

	if err != nil {

		//	formattedError := formaterror.FormatError(err.Error())

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, hotelCreated)
}
