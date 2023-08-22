-- +goose Up
CREATE TABLE comment_votes (
  id SERIAL NOT NULL PRIMARY KEY,
  comment_id INT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type vote_type NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

  UNIQUE(user_id, comment_id)
);

-- +goose Down
DROP TABLE comment_votes;
