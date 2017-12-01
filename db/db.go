package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
)

var db *sqlx.DB

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
	Buyin     int       `db:"buyin"`
	Walkout   int       `db:"walkout"`
	CreatedAt time.Time `db:"created_at"`
}

// InitDB initializes the database connection
func InitDB(dbType string, dbInfo string) {
	var err error
	db, err = sqlx.Open(dbType, dbInfo)

	if err != nil {
		panic(err)
	}
}

// Close database method
func Close() {
	db.Close()
}

// CreateRealm method
func CreateRealm(name null.String, title null.String) (int, error) {
	insertRealm := `
		INSERT INTO realm (name, title)
		VALUES ($1, $2)
		RETURNING id
	`

	var realmID int
	err := db.QueryRow(insertRealm, name, title).Scan(&realmID)

	if err != nil {
		fmt.Println("Failed to create new realm:", err.Error())
		return 0, err
	}

	fmt.Printf("Successfully created new realm id=%d name=\"%s\"", realmID, name.String)
	return realmID, nil
}

// CreatePlayer method
func CreatePlayer(name null.String, realmID null.Int) (int, error) {
	insertPlayer := `
		INSERT INTO player (name, realm_id)
		VALUES ($1, $2)
		RETURNING id
	`

	var playerID int
	err := db.QueryRow(insertPlayer, name, realmID).Scan(&playerID)

	if err != nil {
		fmt.Println("Failed to create new player")
		return 0, err
	}

	fmt.Printf("Successfully created new player id=%d name=\"%s\"\n", playerID, name.String)
	return playerID, nil
}
