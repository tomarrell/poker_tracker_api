# Poker API

## Vision for this project

The goal for this project is for a user to be able to come to the site, enter a realm and be presented with a leaderboard, which shows an overview of all the players for each session up to the latest, showing interesting stats such as how much money has changed hands, the richest players etc. Eye candy graphs are a bonus. 

They can then browse through the sessions sorted by the time they occurred. Viewing a session will show the players who were present, their winnings and losses during that session, and the total that they were on up to that point.

The user is then able to create a new session, this will present him with a list of the regular player which he can tick to indicate they played and enter in their buyin/walkout amounts. At the bottom of the list will be a button to create a new player, this will add the player to the `player` table and link them to that realm. Once the new session has been filled out and submitted, a record of each player who played will be inserted into the mapping table `player_session` tying the player to a newly created session in the `session` table.

## Prod DB

Currently running in heroku as an add on. Ask one of the contributors for prod credentials.

## Prod Deployment

Currently deployed on heroku, running on [https://poker-tracker-api.herokuapp.com/](https://poker-tracker-api.herokuapp.com/)

To deploy to heroku you'll need the heroku CLI installed. Then...

```
heroku login
heroku git:remote -a poker-tracker-api
git push heroku [your-branch-name-here]:master
```

## Prerequisites

- golang version 1.8 or greater
- https://github.com/golang/dep
- docker && docker-compose

## Setup

1. Clone the repository
2. `dep ensure` to resolve dependencies
3. Run `docker-compose up` to bring up local postgres db
4. Install https://github.com/pressly/goose and run Goose migration command described below.
5. Run `go install && poker_tracker_api --config=config.yaml` inside the root directory to run the service on port `8080`
6. Navigate to `localhost:8080/` to view graphql interface

### Goose migration command
For local after running docker-compose:
```
goose --dir migrations postgres "host=localhost user=postgres dbname=pokerapi password=redsux sslmode=disable" up
```

### Tests

Running all tests:
```
docker-compose up -d
go test -v
```

Run short without integration tests:
```
go test -v --short
```


## Todo:
See active [Issues](https://github.com/tomarrell/poker_tracker_api/issues) 

