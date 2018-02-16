-- +goose Up
ALTER TABLE "player"
ADD CONSTRAINT player_name_realm_id UNIQUE(name, realm_id);

-- +goose Down
ALTER TABLE "player"
DROP CONSTRAINT player_name_realm_id;
