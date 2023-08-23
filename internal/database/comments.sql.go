// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: comments.sql

package database

import (
	"context"
)

const findCommentsOfAPost = `-- name: FindCommentsOfAPost :many
SELECT id, text, post_id, author_id, reply_to_id, created_at, updated_at FROM comments WHERE post_id = $1
`

func (q *Queries) FindCommentsOfAPost(ctx context.Context, postID int32) ([]Comment, error) {
	rows, err := q.db.QueryContext(ctx, findCommentsOfAPost, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Comment
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.ID,
			&i.Text,
			&i.PostID,
			&i.AuthorID,
			&i.ReplyToID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
