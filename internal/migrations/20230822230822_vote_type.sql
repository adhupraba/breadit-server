-- +goose Up
CREATE TYPE vote_type AS ENUM ('UP', 'DOWN');

-- +goose Down
DROP TYPE vote_type;
