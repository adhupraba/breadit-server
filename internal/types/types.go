package types

import "github.com/adhupraba/breadit-server/internal/database"

type Json map[string]any

type PostWithAuthorAndVotes struct {
	database.Post
	Author database.User   `json:"author"`
	Votes  []database.Vote `json:"votes"`
}

type CommentWithAuthorAndVotes struct {
	database.Comment
	Author database.User          `json:"author"`
	Votes  []database.CommentVote `json:"votes"`
}

type CommentWithReplies struct {
	CommentWithAuthorAndVotes
	Replies []CommentWithAuthorAndVotes `json:"replies"`
}
