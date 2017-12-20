-- +goose Up
CREATE TABLE "transfer" (
    "id" serial,
    "player_id" int NOT NULL,
    "amount" int NOT NULL,
    "session_id" int,
    "reason" text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("player_id") REFERENCES "player"("id") ON DELETE CASCADE,
    FOREIGN KEY ("session_id") REFERENCES "session"("id") ON DELETE CASCADE
);

-- +goose Down
DROP TABLE "transfer"