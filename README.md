# Poker API

## Vision for this project

The goal for this project is for a user to be able to come to the site, enter a realm and be presented with a leaderboard, which shows an overview of all the players for each session up to the latest, showing interesting stats such as how much money has changed hands, the richest players etc. Eye candy graphs are a bonus. 

They can then browse through the sessions sorted by the time they occurred. Viewing a session will show the players who were present, their winnings and losses during that session, and the total that they were on up to that point.

The user is then able to create a new session, this will present him with a list of the regular player which he can tick to indicate they played and enter in their buyin/walkout amounts. At the bottom of the list will be a button to create a new player, this will add the player to the `player` table and link them to that realm. Once the new session has been filled out and submitted, a record of each player who played will be inserted into the mapping table `player_session` tying the player to a newly created session in the `session` table.

## Prod DB

I currently have a DB running on Google Cloud Platform, it requires a whitelisted IP however and will be used as the prod DB. The details are also included below and in the code as I was using it during development when experimenting with GCP. 

## Prerequisites

- golang version 1.8 or greater
- https://github.com/golang/dep
- docker && docker-compose

## Setup

1. Clone the repository
2. `dep ensure` to resolve dependencies
3. Run `docker-compose up` to bring up local postgres db
4. `cd` into `/migrations` and with https://github.com/pressly/goose run the goose command below, with your DB details inserted to run the migrations and setup your database
5. Run `go install && poker_tracker_api --config=config.yaml` inside the root directory to run the service on port `8080`
6. Navigate to `localhost:8080/` to view graphql interface

### Goose migration command
For GCP: `goose postgres "host=35.197.168.240 user=postgres dbname=postgres password=gl1iKw8B1OCPIM5A sslmode=disable" up`

For local after running docker-compose: `goose postgres "host=localhost user=postgres dbname=postgres password=crimsonsux sslmode=disable" up`

## Todo:
See active [Issues](https://github.com/tomarrell/poker_tracker_api/issues) 

