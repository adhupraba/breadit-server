-- +goose Up
CREATE TABLE subscriptions (
  id SERIAL NOT NULL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  subreddit_id INT NOT NULL REFERENCES subreddits(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

  UNIQUE(user_id, subreddit_id)
);

-- +goose Down
DROP TABLE subscriptions;
