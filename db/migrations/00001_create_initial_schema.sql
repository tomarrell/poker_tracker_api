-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE "realm" (
    "id" serial,
    "name" text NOT NULL UNIQUE,
    "title" text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE "session" (
    "id" serial,
    "realm_id" int NOT NULL,
    "name" text,
    "time" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("realm_id") REFERENCES "realm"("id") ON DELETE CASCADE
);

CREATE TABLE "player" (
    "id" serial,
    "name" text NOT NULL,
    "realm_id" int NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("realm_id") REFERENCES "realm"("id") ON DELETE CASCADE
);

CREATE TABLE "player_session" (
    "player_id" int,
    "session_id" int,
    "buyin" int NOT NULL,
    "walkout" int NOT NULL,
    PRIMARY KEY ("player_id", "session_id"),
    FOREIGN KEY ("player_id") REFERENCES "player"("id") ON DELETE CASCADE,
    FOREIGN KEY ("session_id") REFERENCES "session"("id") ON DELETE CASCADE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE "realm";
DROP TABLE "session";
DROP TABLE "player";
DROP TABLE "player_session";
