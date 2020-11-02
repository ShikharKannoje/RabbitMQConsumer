package controller

import (
	"RabbitMQConsumer/consumer/model"
	"RabbitMQConsumer/consumer/response"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (server *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	room := model.Room{}
	err = json.Unmarshal(body, &room)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	room.Prepare()
	err = room.Validate()
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	roomCreated, err := room.SaveRoom(server.DB)

	if err != nil {

		//	formattedError := formaterror.FormatError(err.Error())

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, roomCreated)
}
