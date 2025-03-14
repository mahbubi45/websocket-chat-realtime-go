package controller

func (s *Server) Routes() {
	// r := s.R.PathPrefix("/api").Subrouter()
	r := s.R

	r.HandleFunc("/selected-user-group", s.UpdateSelectedUserToGrupController).Methods("PUT")
	r.HandleFunc("/add-group", s.AddGrupMemberUsersController).Methods("POST")

	r.HandleFunc("/ws/group", s.HandleConnectionsGrupController).Methods("GET")
	r.HandleFunc("/ws/private", s.HandleConnectionsPrivateMessageController).Methods("GET")
}
