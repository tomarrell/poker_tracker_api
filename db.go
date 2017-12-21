package main

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/guregu/null.v3"
)

// Realm DB table
type Realm struct {
	ID        int         `db:"id"`
	Name      string      `db:"name"`
	Title     null.String `db:"title"`
	CreatedAt time.Time   `db:"created_at"`
}

// Session DB table
type Session struct {
	ID        int         `db:"id"`
	RealmID   int         `db:"realm_id"`
	Name      null.String `db:"name"`
	Time      time.Time   `db:"time"`
	CreatedAt time.Time   `db:"created_at"`
}

// Player DB table
type Player struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	RealmID   int       `db:"realm_id"`
	CreatedAt time.Time `db:"created_at"`
}

// PlayerSession DB table
type PlayerSession struct {
	PlayerID  int       `db:"player_id"`
	SessionID int       `db:"session_id"`
	Buyin     null.Int  `db:"buyin"`
	Walkout   null.Int  `db:"walkout"`
	CreatedAt time.Time `db:"created_at"`
}

type postgresDb struct {
	db *sqlx.DB
}

// InitDB initializes the database connection
func mustInitDB(dbType string, dbInfo string) *postgresDb {

	db := sqlx.MustOpen(dbType, dbInfo)
	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &postgresDb{db}
}

// Close database method
func (p *postgresDb) Close() {
	p.db.Close()
}

// CreateRealm method
func (p *postgresDb) CreateRealm(name null.String, title null.String) (*Realm, error) {
	insertRealm := `
		INSERT INTO realm (name, title)
		VALUES ($1, $2)
		RETURNING id, name, title, created_at
	`

	var realm Realm
	if err := p.db.Get(&realm, insertRealm, name, title); err != nil {
		log.WithError(err).Error("Failed to create new realm:", err.Error())
		return nil, err
	}

	log.Infof(`Successfully created new realm id=%d name="%s"`, realm.ID, name.String)
	return &realm, nil
}

// CreatePlayer method
func (p *postgresDb) CreatePlayer(name null.String, realmID null.Int) (*Player, error) {
	insertPlayer := `
		INSERT INTO player (name, realm_id)
		VALUES ($1, $2)
		RETURNING id, name, realm_id
	`
	var player Player
	err := p.db.Get(&player, insertPlayer, name, realmID)

	if err != nil {
		log.WithError(err).Error("Failed to create new player")
		return nil, err
	}

	fmt.Printf("Successfully created new player id=%d name=\"%s\"\n", player.ID, name.String)
	return &player, nil
}

// CreateSession handles creating the session row and creating player_session records
func (p *postgresDb) CreateSession(realmID null.Int, name null.String, time null.Time, playerSessions []PlayerSession) (*Session, error) {
	insertSession := `
		INSERT INTO session (realm_id, name, time)
		VALUES ($1, $2, $3)
		RETURNING id, realm_id, name, time
	`
	mapPlayerToSession := `
		INSERT INTO player_session (player_id, session_id, buyin, walkout)
		VALUES ($1, $2, $3, $4)
	`

	insertTransferForPlayer := `
		INSERT INTO transfer (player_id, amount, session_id, reason)
		VALUES($1, $2, $3, $4)
	`

	var session Session

	// Begin transaction
	tx := p.db.MustBegin()
	if err := tx.Get(&session, insertSession, realmID, name, time); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, player := range playerSessions {
		if _, err := tx.Exec(mapPlayerToSession, player.PlayerID, session.ID, player.Buyin, player.Walkout); err != nil {
			tx.Rollback()
			return nil, err
		}
		var amount = player.Walkout.Int64 - player.Buyin.Int64
		if _, err := tx.Exec(insertTransferForPlayer, player.PlayerID, amount, session.ID, "Session Participation"); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err := tx.Commit()

	if err != nil {
		log.WithError(err).Error("Failed to record new session")
		return nil, err
	}

	log.Infof(`Successfully created new session id=%d name="%s" time="%s"`, session.ID, session.Name.String, session.Time.String())
	return &session, nil
}

// GetSessions returns an array of sessions given a realmID
func (p *postgresDb) GetSessions(realmID int) ([]Session, error) {
	getSessions := `
		SELECT *
		FROM session
		WHERE realm_id=$1
	`

	var sessions []Session
	if err := p.db.Select(&sessions, getSessions, realmID); err != nil {
		log.WithError(err).Error("Failed to fetch sessions of realm")
		return nil, err
	}

	return sessions, nil
}

func (p *postgresDb) GetSessionByID(id int) (*Session, error) {
	getSessions := `
		SELECT *
		FROM session
		WHERE id=$1
	`

	var session Session
	if err := p.db.Get(&session, getSessions, id); err != nil {
		log.WithError(err).Error("Failed to fetch sessions by ID")
		return nil, err
	}

	return &session, nil
}

func (p *postgresDb) GetPlayerByID(id int) (*Player, error) {
	getPlayers := `
		SELECT *
		FROM player
		WHERE id=$1
	`

	var player Player
	if err := p.db.Get(&player, getPlayers, id); err != nil {
		log.WithError(err).Error("Failed to fetch player by ID")
		return nil, err
	}

	return &player, nil
}

func (p *postgresDb) GetRealmByName(name string) (*Realm, error) {
	q := `
		SELECT *
		FROM realm
		WHERE name=$1
	`

	var realm Realm
	if err := p.db.Get(&realm, q, name); err != nil {
		log.WithError(err).Error("Failed to fetch realm by name")
		return nil, err
	}

	return &realm, nil
}
