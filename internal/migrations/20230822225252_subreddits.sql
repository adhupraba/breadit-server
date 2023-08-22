-- +goose Up
-- +goose StatementBegin
CREATE TABLE subreddits (
  id SERIAL NOT NULL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  creator_id INT REFERENCES users(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX subreddits_name_idx ON subreddits("name");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX subreddits_name_idx;

DROP TABLE subreddits;
-- +goose StatementEnd
