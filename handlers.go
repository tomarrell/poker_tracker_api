package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tomarrell/poker_api/db"
)

// Realm Management

// CreateRealmHandler method
// Params--
// 	Name:  string
// 	Title: ?string
func CreateRealmHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var realm db.Realm

	err := decoder.Decode(&realm)
	if err != nil {
		RespondServerError(w, []string{
			"Failed to decode JSON",
			err.Error(),
		})
		return
	}

	realmID, err := db.CreateRealm(realm.Name, realm.Title)

	if err != nil {
		RespondServerError(w, []string{
			"Failed to create realm",
			err.Error(),
		})
		return
	}

	resp := struct{ RealmID int }{realmID}
	fmt.Println(realmID)
	SuccessWithJSON(w, resp)
}

// Player management

// CreatePlayerHandler method
// Params--
//   Name:    string
//   RealmID: string
func CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var player db.Player

	err := decoder.Decode(&player)
	if err != nil {
		RespondServerError(w, []string{
			"Failed to decode JSON",
			err.Error(),
		})
		return
	}

	playerID, err := db.CreatePlayer(player.Name, player.RealmID)

	if err != nil {
		RespondServerError(w, []string{
			"Failed to create new player",
			err.Error(),
		})
		return
	}

	resp := struct{ PlayerID int }{playerID}
	fmt.Println(playerID)
	SuccessWithJSON(w, resp)
}

// Session management

// CreateSessionHandler method
func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating new session...")
}
