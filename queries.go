package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type queryHandler struct {
	db *postgresDb
}

func (q *queryHandler) GetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars["realmID"])

	realmID, err := strconv.Atoi(vars["realmID"])

	if err != nil {
		RespondServerError(w, []string{"Failed to decode realmID from params", err.Error()})
		return
	}

	sessions, err := q.db.GetSessions(realmID)

	if err != nil {
		RespondServerError(w, []string{"Failed to fetch sessions of realm", err.Error()})
		return
	}

	resp := struct{ Sessions []Session }{sessions}
	SuccessWithJSON(w, resp)
}
