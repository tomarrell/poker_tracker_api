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
	ID      int    `db:"id"`
	Name    string `db:"name"`
	RealmID int    `db:"realm_id"`
}

// PlayerSession DB table
type PlayerSession struct {
	PlayerID  int `db:"player_id"`
	SessionID int `db:"session_id"`
	Buyin     int `db:"buyin"`
	Walkout   int `db:"walkout"`
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
func CreateRealm(name string, title null.String) error {
	insertRealm := `INSERT INTO realm (name, title) VALUES ($1, $2)`
	_, err := db.Exec(insertRealm, name, title)

	if err != nil {
		fmt.Println("Failed to create new Realm")
		return err
	}

	fmt.Println("Successfully created new realm:", name)
	return nil
}
