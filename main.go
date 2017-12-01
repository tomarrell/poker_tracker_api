package main

import (
	"fmt"
	"net/http"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/gorilla/mux"
	"github.com/tomarrell/poker_tracker_api/db"
)

func main() {
	dbInfo := `
		host=tom-personal:australia-southeast1:tom-personal
		user=postgres
		dbname=postgres
		password=gl1iKw8B1OCPIM5A
		sslmode=disable
	`
	dbType := "cloudsqlpostgres"
	PORT := ":3000"

	db.InitDB(dbType, dbInfo)
	defer db.Close()

	// Set up routing
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/new/realm", CreateRealmHandler).Methods("POST")
	api.HandleFunc("/new/player", CreatePlayerHandler).Methods("POST")
	api.HandleFunc("/new/session", CreateSessionHandler).Methods("POST")

	fmt.Println("HTTP Server opening on port", PORT)
	http.ListenAndServe(PORT, r)
}
