package rawmessageparser

import (
	"encoding/json"
	"time"

	"github.com/adhupraba/breadit-server/internal/database"
)

func ParseJsonVote(data json.RawMessage) (database.Vote, error) {
	var parsed map[string]interface{}

	err := json.Unmarshal(data, &parsed)

	if err != nil {
		return database.Vote{}, err
	}

	dbVote, err := transformJsonVote(parsed)

	return dbVote, err
}

func ParseJsonVotes(data json.RawMessage) ([]database.Vote, error) {
	var parsed []map[string]interface{}

	err := json.Unmarshal(data, &parsed)

	if err != nil {
		return []database.Vote{}, err
	}

	votes := []database.Vote{}

	for _, vote := range parsed {
		if vote == nil || int32(vote["id"].(float64)) == 0 {
			continue
		}

		dbVote, err := transformJsonVote(vote)

		if err != nil {
			return []database.Vote{}, err
		}

		votes = append(votes, dbVote)
	}

	return votes, nil
}

func transformJsonVote(voteMap map[string]interface{}) (database.Vote, error) {
	createdAt, err := time.Parse(time.RFC3339, voteMap["created_at"].(string)+"Z")
	updatedAt, err := time.Parse(time.RFC3339, voteMap["updated_at"].(string)+"Z")

	if err != nil {
		return database.Vote{}, err
	}

	vote := database.Vote{
		ID:        int32(voteMap["id"].(float64)),
		PostID:    int32(voteMap["post_id"].(float64)),
		UserID:    int32(voteMap["user_id"].(float64)),
		Type:      database.VoteType(voteMap["type"].(string)),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return vote, nil
}
