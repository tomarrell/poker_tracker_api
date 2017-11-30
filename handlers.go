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

	realmID, err := db.CreateRealm(realm.Name, realm.Title)
	if err != nil {
		RespondServerError(
			w,
			[]string{
				"Failed to create realm.",
				err.Error(),
			},
		)
		return
	}

	resp := struct{ ID int }{realmID}
	fmt.Println(realmID)
	SuccessWithJSON(w, resp)
}

// Player management

// CreatePlayerHandler method
func CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
}

// Session management

// CreateSessionHandler method
func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating new session...")
}
