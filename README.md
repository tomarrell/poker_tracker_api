# Poker API

I currently have a DB running on Google Cloud Platform, it requires a whitelisted IP however and will be used as the prod DB. The details are also included below and in the code as I was using it during development when experimenting with GCP. 

## Setup
1. Clone the repository
2. `go get` to fetch Go dependencies
3. Run a local version of postgres, change the hostname, user, dbname and password info in main.go
4. `cd` into `/db/migrations` and with https://github.com/pressly/goose run the goose command below, with your DB details inserted to run the migrations and setup your database
5. Run `go run` inside the root directory to run the service on port `3000`

### Goose migration command
For GCP: `goose postgres "host=35.197.168.240 user=postgres dbname=postgres password=gl1iKw8B1OCPIM5A sslmode=disable" up`

For local: `goose postgres "host=[HOST] user=[USER] dbname=[DB_NAME] password=[PASSWORD] sslmode=disable" up`

## Todo:
- Implement better request body validation
- Cleanup response methods
