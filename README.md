# Poker API

goose postgres "host=35.197.168.240 user=postgres dbname=postgres password=gl1iKw8B1OCPIM5A sslmode=disable" up

## Setup
1. Clone the repository
2. `go get` to fetch Go dependencies
3. Run a local version of postgres, change the hostname, user, dbname and password info in main.go
4. `cd` into `/db/migrations` and run with https://github.com/pressly/goose run the goose command above, with your DB details inserted to run the migrations and setup your database
5. Run `go run` inside the root directory to run the service on port `3000`

## Todo:
- Implement better request body validation
- Cleanup response methods
