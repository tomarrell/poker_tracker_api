package main

import (
	"errors"
	"time"

	"gopkg.in/guregu/null.v3"
)

// CreateRealm args struct
type CreateRealm struct {
	Name  string
	Title *string
}

// CreateRealm mutation
func (r *Resolver) CreateRealm(args CreateRealm) (*RealmResolver, error) {
	if args.Name == "" {
		return nil, errors.New("Must supply realm name")
	}

	realm, err := r.db.CreateRealm(args.Name, args.Title)
	if err != nil {
		return nil, err
	}

	return &RealmResolver{
		id:    toGQL(realm.ID),
		name:  realm.Name,
		title: realm.Title.Ptr(),
		db:    r.db}, nil
}

// CreatePlayer args struct
type CreatePlayer struct {
	Name    string
	RealmID int32
}

// CreatePlayer mutation
func (r *Resolver) CreatePlayer(args CreatePlayer) (*PlayerResolver, error) {
	if args.Name == "" {
		return nil, errors.New("Must supply player name")
	}
	if args.RealmID == 0 {
		return nil, errors.New("Must supply Realm ID")
	}

	// TODO: Check if realm exists first to provide nicer error?
	player, err := r.db.CreatePlayer(args.Name, args.RealmID)

	if err != nil {
		return nil, err
	}

	return &PlayerResolver{
		id:      toGQL(player.ID),
		name:    player.Name,
		realmID: toGQL(player.RealmID),
	}, nil
}

// CreateSessionPlayerSession args struct
type CreateSessionPlayerSession struct {
	PlayerID int32
	Buyin    int32
	Walkout  int32
}

// CreateSession args struct
type CreateSession struct {
	Name           string
	RealmID        int32
	Time           string
	PlayerSessions []*CreateSessionPlayerSession
}

// CreateSession mutation
func (r *Resolver) CreateSession(args CreateSession) (*SessionResolver, error) {
	// TODO: Argument validation

	var t *time.Time
	if args.Time != "" {
		parsedTime, err := time.Parse(time.RFC3339, args.Time)
		if err != nil {
			return nil, err
		}
		t = &parsedTime
	}

	var playerSessions = make([]PlayerSession, len(args.PlayerSessions))
	for i, csps := range args.PlayerSessions {
		playerSessions[i] = PlayerSession{
			PlayerID: int(csps.PlayerID),
			Buyin:    null.IntFrom(int64(csps.Buyin)),
			Walkout:  null.IntFrom(int64(csps.Walkout)),
		}
	}

	session, err := r.db.CreateSession(
		args.RealmID,
		args.Name,
		t,
		playerSessions)

	if err != nil {
		return nil, err
	}
	return &SessionResolver{
		id:      toGQL(session.ID),
		name:    session.Name.Ptr(),
		realmID: toGQL(session.RealmID),
		time:    session.Time.Format(time.RFC3339),
	}, nil
}
