package controller

import (
	"RabbitMQConsumer/consumer/response"
	"net/http"
)

func (server *Server) ome(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "WI")

}
