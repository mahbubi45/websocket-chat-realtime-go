package controller

func (s *Server) Routes() {
	// r := s.R.PathPrefix("/api").Subrouter()
	r := s.R

	r.HandleFunc("/ws/group", s.HandleConnectionsGrupController).Methods("GET")
	r.HandleFunc("/ws/private", s.HandleConnectionsPrivateMessageController).Methods("GET")
}
