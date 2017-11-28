package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tomarrell/poker_api/db"
)

// Realm Management

// CreateRealmHandler method
// Available params:
// 	name:string
// 	title:?string
func CreateRealmHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var realm db.Realm

	err := decoder.Decode(&realm)
	if err != nil {
		panic(err)
	}

	err = db.CreateRealm(realm.Name, realm.Title)
	if err != nil {
		http.Error(w, "Failed to create realm.\n\t"+err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "OK")
}

// Player management

// CreatePlayerHandler method
func CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating new player...")
}

// Session management

// CreateSessionHandler method
func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating new session...")
}
