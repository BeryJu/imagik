package server

import "net/http"

func (s *Server) HealthLiveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
}

func (s *Server) HealthReadiness(w http.ResponseWriter, r *http.Request) {
	if s.HashMap.Populated() {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(500)
	}
}
