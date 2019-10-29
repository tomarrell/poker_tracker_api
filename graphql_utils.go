package main

import (
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
)

func toGQL(i int) graphql.ID {
	return graphql.ID(strconv.Itoa(i))
}
