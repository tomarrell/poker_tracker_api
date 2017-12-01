package main

import (
	"encoding/json"
	"net/http"

	"github.com/tomarrell/poker_tracker_api/db"
	"gopkg.in/guregu/null.v3"
)

// Realm Management

type realmRequest struct {
	Name  null.String
	Title null.String
}

// CreateRealmHandler handles realm creation
func CreateRealmHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var realm realmRequest

	err := decoder.Decode(&realm)

	if err != nil {
		FailDecode(w, err)
		return
	}

	if realm.Name.IsZero() {
		InvalidArgs(w, []string{"Name"})
		return
	}

	realmID, err := db.CreateRealm(realm.Name, realm.Title)

	if err != nil {
		RespondServerError(w, []string{"Failed to create realm", err.Error()})
		return
	}

	resp := struct{ RealmID int }{realmID}
	SuccessWithJSON(w, resp)
}

// Player management

type playerRequest struct {
	Name    null.String
	RealmID null.Int
}

// CreatePlayerHandler handles player creation
func CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var player playerRequest

	err := decoder.Decode(&player)

	if err != nil {
		FailDecode(w, err)
		return
	}

	if player.Name.IsZero() {
		InvalidArgs(w, []string{"Name"})
		return
	}

	if player.RealmID.IsZero() {
		InvalidArgs(w, []string{"RealmID"})
		return
	}

	playerID, err := db.CreatePlayer(player.Name, player.RealmID)

	if err != nil {
		RespondServerError(w, []string{"Failed to create new player", err.Error()})
		return
	}

	resp := struct{ PlayerID int }{playerID}
	SuccessWithJSON(w, resp)
}

// Session management

type sessionRequest struct {
	RealmID   null.Int
	Name      null.String
	Time      null.Time
	PlayerIDs []int
}

// CreateSessionHandler handles session creation
func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var sess sessionRequest

	err := decoder.Decode(&sess)

	if err != nil {
		FailDecode(w, err)
	}

	if sess.RealmID.IsZero() {
		InvalidArgs(w, []string{"RealmID"})
		return
	}

	if !sess.Time.Valid {
		InvalidArgs(w, []string{"Time"})
		return
	}

	if len(sess.PlayerIDs) < 1 {
		InvalidArgs(w, []string{"PlayerIDs"})
		return
	}

	sessionID, err := db.CreateSession(sess.RealmID, sess.Name, sess.Time, sess.PlayerIDs)

	if err != nil {
		RespondServerError(w, []string{"Failed to create new session", err.Error()})
		return
	}

	resp := struct{ SessionID int }{sessionID}
	SuccessWithJSON(w, resp)
}
