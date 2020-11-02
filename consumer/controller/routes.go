package controller

func (s *Server) InitializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", s.ome).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/app", s.app).Methods("POST")

}
