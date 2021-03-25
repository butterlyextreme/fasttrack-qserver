// server.go

package main

import (
	"log"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func (s *Server) Initialize() {
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", s.Router))
}

func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/questions", s.getQuestions).Methods("GET")
	s.Router.HandleFunc("/result", s.calculateResults).Methods("POST")
}

func (s *Server) getQuestions(w http.ResponseWriter, r *http.Request) {
	answers, err := getQuestions()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, answers)
}

func (s *Server) calculateResults(w http.ResponseWriter, r *http.Request) {
	var ans answers
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ans); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	score, err := ans.calculateResults()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, score)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
