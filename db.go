package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
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

// Transfer DB table
type Transfer struct {
	ID        int         `db:"id"`
	PlayerID  int         `db:"player_id"`
	Amount    int         `db:"amount"`
	SessionID null.Int    `db:"session_id"`
	Reason    null.String `db:"reason"`
	CreatedAt time.Time   `db:"created_at"`
}

type balanceSummary struct {
	PlayerName string `db:"name"`
	Total      int    `db:"total"`
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
func (p *postgresDb) CreateRealm(name string, title *string) (*Realm, error) {
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

	log.Debugf(`Successfully created new realm id=[%d] name=[%s]`, realm.ID, name)
	return &realm, nil
}

// CreatePlayer method
func (p *postgresDb) CreatePlayer(name string, realmID int32) (*Player, error) {
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

	log.Debugf(`Successfully created new player id=[%d] name=[%s]`, player.ID, name)
	return &player, nil
}

// CreateSession handles creating or updating the session row and player_session records
// It follows put semantics

func (p *postgresDb) CreateOrUpdateSession(sessionID *int32, realmID int32, name string, t *time.Time, playerSessions []PlayerSession) (*Session, error) {

	deletePlayerSessions := `
	DELETE FROM player_session WHERE session_id = $1
	`
	deletePlayerTransfers := `
	DELETE FROM transfer WHERE session_id = $1
	`
	updateSession := `
	UPDATE session SET realm_id=$1, name=$2, time=$3 WHERE id=$4
	RETURNING id, realm_id, name, time
	`

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

	if sessionID != nil {
		if _, err := tx.Exec(deletePlayerSessions, *sessionID); err != nil {
			tx.Rollback()
			return nil, err
		}
		if _, err := tx.Exec(deletePlayerTransfers, *sessionID); err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Get(&session, updateSession, realmID, name, t, *sessionID); err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		if err := tx.Get(&session, insertSession, realmID, name, t); err != nil {
			tx.Rollback()
			return nil, err
		}
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
		log.WithError(err).Error("Failed to record session upsert")
		return nil, err
	}

	log.Debugf(`Successfully upsert session id=[%d] name=[%s] time=[%s]`, session.ID, session.Name.String, session.Time.String())
	return &session, nil
}

// GetSessions returns an array of sessions given a realmID
func (p *postgresDb) GetSessionsByRealmID(realmID int) ([]Session, error) {
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

func (p *postgresDb) GetPlayersByRealmID(realmID int) ([]Player, error) {
	getPlayers := `
		SELECT *
		FROM player
		WHERE realm_id=$1
	`

	var players []Player
	if err := p.db.Select(&players, getPlayers, realmID); err != nil {
		log.WithError(err).Error("Failed to fetch players by realm ID")
		return nil, err
	}

	return players, nil
}

// GetRealmByField should only be used with fields that are unique ids, i.e. name and id
func (p *postgresDb) GetRealmByField(field string, val interface{}) (*Realm, error) {
	q := fmt.Sprintf(`
		SELECT *
		FROM realm
		WHERE %s=$1
	`, field)

	var realm Realm
	if err := p.db.Get(&realm, q, val); err != nil {
		log.WithError(err).Errorf("Failed to fetch realm by %s", field)
		return nil, err
	}

	return &realm, nil
}

func (p *postgresDb) GetPlayerSessionsByField(field string, val interface{}) ([]PlayerSession, error) {
	q := fmt.Sprintf(`
		SELECT *
		FROM player_session
		WHERE %s=$1
	`, field)

	var pSessions []PlayerSession
	if err := p.db.Select(&pSessions, q, val); err != nil {
		log.WithError(err).Errorf("Failed to fetch player_sessions by %s", field)
		return nil, err
	}

	return pSessions, nil
}

func (p *postgresDb) GetHistoricalBalanceByPlayerID(id int) (int32, error) {
	q := `
		SELECT COALESCE(SUM(walkout) - SUM(buyin), 0)
		FROM player_session
		WHERE player_id=$1
	`

	var balance int32
	if err := p.db.Get(&balance, q, id); err != nil {
		log.WithError(err).Errorf("Failed to get player session stats for player id %v", id)
		return 0, err
	}
	return balance, nil
}

func (p *postgresDb) GetRealBalanceByPlayerID(id int) (int32, error) {
	q := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transfer
		WHERE player_id=$1
	`

	var balance int32
	if err := p.db.Get(&balance, q, id); err != nil {
		log.WithError(err).Errorf("Failed to get transfer stats for player id %v", id)
		return 0, err
	}
	return balance, nil
}

func (p *postgresDb) GetBalanceSummaryByPlayerIDs(ids []int) ([]balanceSummary, error) {
	q := `
		SELECT p.name as name, COALESCE(SUM(amount), 0) as total
		FROM transfer t
		INNER JOIN player p ON t.player_id = p.id
		WHERE t.player_id IN (?)
		GROUP BY name
		ORDER BY total DESC
	`
	query, args, err := sqlx.In(q, ids)
	if err != nil {
		log.WithError(err)
	}
	query = p.db.Rebind(query)

	var balanceSummaries []balanceSummary

	if err := p.db.Select(&balanceSummaries, query, args...); err != nil {
		log.WithError(err).Errorf("Failed to get transfer summaries for %v", ids)
		return nil, err
	}
	return balanceSummaries, nil
}

func (p *postgresDb) GetTotalBuyinByPlayerID(id int) (int32, error) {
	q := `
		SELECT COALESCE(SUM(buyin), 0)
		FROM player_session
		WHERE player_id=$1
	`

	var balance int32
	if err := p.db.Get(&balance, q, id); err != nil {
		log.WithError(err).Errorf("Failed to get total buyin for player id %v", id)
		return 0, err
	}
	return balance, nil
}
