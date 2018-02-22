package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"
)

// CreateRealm args struct
type CreateRealm struct {
	Name  string
	Title *string
}

type byNet []PlayerSession

func (n byNet) Len() int {
	return len(n)
}

func (n byNet) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n byNet) Less(i, j int) bool {
	return (n[i].Walkout.Int64 - n[i].Buyin.Int64) > (n[j].Walkout.Int64 - n[j].Buyin.Int64)
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
	playerName := strings.ToLower(strings.Trim(args.Name, " "))

	if playerName == "" {
		return nil, errors.New("Must supply player name")
	}
	if args.RealmID == 0 {
		return nil, errors.New("Must supply Realm ID")
	}

	// TODO: Check if realm exists first to provide nicer error?
	player, err := r.db.CreatePlayer(playerName, args.RealmID)

	if err != nil {
		return nil, err
	}

	return &PlayerResolver{
		id:      toGQL(player.ID),
		name:    player.Name,
		realmID: toGQL(player.RealmID),
		db:      r.db,
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
	ID             *int32
}

// CreateSession mutation
func (r *Resolver) PutSession(args CreateSession) (*SessionResolver, error) {

	if args.Name == "" {
		return nil, errors.New("Session name cannot be empty")
	}

	if args.RealmID == 0 {
		return nil, errors.New("Realm id cannot be empty or 0")
	}

	var t *time.Time
	if args.Time != "" {
		parsedTime, err := time.Parse(time.RFC3339, args.Time)
		if err != nil {
			return nil, err
		}
		t = &parsedTime
	} else {
		return nil, errors.New("Time cannot be empty")
	}

	if args.ID != nil {
		if _, err := r.db.GetSessionByID(int(*args.ID)); err != nil {
			return nil, fmt.Errorf("Cannot validate existence of session with id %v", *args.ID)
		}
	}

	if _, err := r.db.GetRealmByField("id", args.RealmID); err != nil {
		return nil, fmt.Errorf("Cannot validate existence of realm with id %v", args.RealmID)
	}

	var (
		playerSessions = make([]PlayerSession, len(args.PlayerSessions))
		pids           = make(map[int32]*Player, len(args.PlayerSessions))
	)
	for i, csps := range args.PlayerSessions {
		if _, ok := pids[csps.PlayerID]; ok {
			return nil, fmt.Errorf("Duplicate player id %v in playerSessions argument", csps.PlayerID)
		}
		player, err := r.db.GetPlayerByID(int(csps.PlayerID))
		if err != nil {
			return nil, fmt.Errorf("Cannot validate existence of player with id %v", csps.PlayerID)
		}
		pids[csps.PlayerID] = player
		playerSessions[i] = PlayerSession{
			PlayerID: int(csps.PlayerID),
			Buyin:    null.IntFrom(int64(csps.Buyin)),
			Walkout:  null.IntFrom(int64(csps.Walkout)),
		}
	}

	session, err := r.db.CreateOrUpdateSession(
		args.ID,
		args.RealmID,
		args.Name,
		t,
		playerSessions)

	if err != nil {
		return nil, err
	}

	spids := make([]int, 0, len(pids))
	pNames := make(map[int]string, len(pids))
	for id, p := range pids {
		spids = append(spids, int(id))
		pNames[int(id)] = p.Name
	}

	bSummary, err := r.db.GetBalanceSummaryByPlayerIDs(spids)
	if err != nil {
		return nil, err
	}
	sort.Sort(byNet(playerSessions))

	r.slacker.SendSummary(bSummary, playerSessions, pNames, int(args.RealmID))

	return &SessionResolver{
		id:      toGQL(session.ID),
		name:    session.Name.Ptr(),
		realmID: toGQL(session.RealmID),
		time:    session.Time.Format(time.RFC3339),
		db:      r.db}, nil
}
