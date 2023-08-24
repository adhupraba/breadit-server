-- +goose Up
CREATE TABLE posts (
  id SERIAL NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  content JSONB,
  subreddit_id INT NOT NULL REFERENCES subreddits(id) ON DELETE CASCADE,
  author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE posts;