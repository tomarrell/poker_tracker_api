package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func main() {
	conf := mustParseConfig()
	db := mustInitDB("postgres", conf.DSN)
	defer db.Close()

	qh := queryHandler{db}
	mh := mutationHandler{db}

	// Set up routing
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/realm", mh.CreateRealmHandler).Methods("POST")
	api.HandleFunc("/player", mh.CreatePlayerHandler).Methods("POST")
	api.HandleFunc("/session", mh.CreateSessionHandler).Methods("POST")

	api.HandleFunc("/sessions/{realmID:[0-9]+}", qh.GetSessionsHandler).Methods("GET")

	log.Infof("starting server on %s", conf.ListenAddress)
	http.ListenAndServe(conf.ListenAddress, r)
}
