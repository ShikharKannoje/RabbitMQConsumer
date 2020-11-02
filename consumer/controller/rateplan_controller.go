package controller

import (
	"RabbitMQConsumer/consumer/model"
	"RabbitMQConsumer/consumer/response"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (server *Server) CreateRateplan(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	rateplan := model.Rateplan{}
	err = json.Unmarshal(body, &rateplan)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	rateplan.Prepare()
	err = rateplan.Validate()
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	rateplanCreated, err := rateplan.SaveRateplan(server.DB)

	if err != nil {

		//	formattedError := formaterror.FormatError(err.Error())

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, rateplanCreated)
}
