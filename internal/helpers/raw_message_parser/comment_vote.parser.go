package rawmessageparser

import (
	"encoding/json"
	"time"

	"github.com/adhupraba/breadit-server/internal/database"
)

func ParseJsonCommentVote(data json.RawMessage) (database.CommentVote, error) {
	var parsed map[string]interface{}

	err := json.Unmarshal(data, &parsed)

	if err != nil {
		return database.CommentVote{}, err
	}

	dbCommentVote, err := transformJsonCommentVote(parsed)

	return dbCommentVote, err
}

func ParseJsonCommentVotes(data json.RawMessage) ([]database.CommentVote, error) {
	var parsed []map[string]interface{}

	err := json.Unmarshal(data, &parsed)

	if err != nil {
		return []database.CommentVote{}, err
	}

	commentVotes := []database.CommentVote{}

	for _, commentVote := range parsed {
		if commentVote == nil || int32(commentVote["id"].(float64)) == 0 {
			continue
		}

		dbCommentVote, err := transformJsonCommentVote(commentVote)

		if err != nil {
			return []database.CommentVote{}, err
		}

		commentVotes = append(commentVotes, dbCommentVote)
	}

	return commentVotes, nil
}

func transformJsonCommentVote(commentVoteMap map[string]interface{}) (database.CommentVote, error) {
	createdAt, err := time.Parse(time.RFC3339, commentVoteMap["created_at"].(string)+"Z")
	updatedAt, err := time.Parse(time.RFC3339, commentVoteMap["updated_at"].(string)+"Z")

	if err != nil {
		return database.CommentVote{}, err
	}

	commentVote := database.CommentVote{
		ID:        int32(commentVoteMap["id"].(float64)),
		CommentID: int32(commentVoteMap["comment_id"].(float64)),
		UserID:    int32(commentVoteMap["user_id"].(float64)),
		Type:      database.VoteType(commentVoteMap["type"].(string)),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return commentVote, nil
}
