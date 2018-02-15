package main

const gqlSchema = `
schema {
		query: Query
		mutation: Mutation
	}

	# The query type, represents all of the entry points into our object graph
	type Query {
		realmByName(name: String!): Realm
		realmById(id: ID!): Realm
		sessionById(id: ID!): Session
		sessionsByRealmId(realmId: ID!): [Session]
		playerById(id: ID!): Player
	}

	# Mutation type, reprsents all modifications available

	type Mutation {
		createRealm(name: String!, title: String): Realm
		createPlayer(name: String!, realmId: ID!): Player
		putSession(id: ID, name: String!, realmId: ID!, time: String!, playerSessions: [CreateSessionPlayerSession]!): Session
	}

	type Player {
		id: ID!
		name: String!
		realmId: ID!
		playerSessions: [PlayerSession]!
		historicalBalance: Int!
		realBalance: Int!
		totalBuyin: Int!
	}

	type Realm {
		id: ID!
		name: String!
		title: String
		players: [Player]!
		sessions: [Session]!
	}

	type Session {
		id: ID!
		realmId: ID!
		name: String
		time: String!
		playerSessions: [PlayerSession]!
	}

	type PlayerSession {
		player: Player!
		playerId: ID!
		sessionId: ID!
		buyin: Int!
		walkout: Int!
	}

	input CreateSessionPlayerSession {
		playerId: ID!
		buyin: Int!
		walkout: Int!
	}
`

// const gqlSchema = `
// schema {
// 		query: Query
// 		mutation: Mutation
// 	}
//
// 	# The query type, represents all of the entry points into our object graph
// 	type Query {
// 		realmByName(name: String!): Realm
// 		realmById(id: ID!): Realm
// 		sessionById(id: ID!): Session
// 		sessionsByRealmId(realmId: ID!): [Session]
// 		playerById(id: ID!): Player
// 		playerSessionByPlayerIdSessionId(playerId: ID!, sessionId: ID!): PlayerSession
// 		playerSessionsByPlayerId(playerId: ID!): [PlayerSession]
// 		playerSessionsBySessionId(sessionId: ID!): [PlayerSession]
// 	}
//
// 	# The mutation type, represents all updates we can make to our data
// 	type Mutation {
// 		createRealm(name: String!, title: String): Realm
// 		createPlayer(name: String!, realmId: ID!): Player
// 		createSession(name: String, realmId: ID!, time: String!, playerSessions: [PlayerSession]!)
// 		updateSession(sessionId: ID!, name: String, time: String, playerSessions: [PlayerSession])
// 	}
//
// 	type Player {
// 		id: ID!
// 		name: String!
// 		realmId: ID!
// 		sessions: [Session]
// 		playerSessions: [PlayerSession]
// 		realm: Realm
// 	}
//
// 	type Realm {
// 		id: ID!
// 		name: String!
// 		title: String
// 		sessions: [Session]
// 		players: [Player]
// 		playerSessions: [PlayerSession]
// 	}
//
// 	type Session {
// 		id: ID!
// 		realmId: ID!
// 		name: String
// 		time: String!
// 		players: [Player]
// 		playerSessions: [PlayerSession]
// 	}
//
// 	type PlayerSession {
// 		playerId: ID!
// 		sessionId: ID!
// 		buyin: Int!
// 		walkout: Int!
// 		player: Player
// 		session: Session
// 	}
// `
