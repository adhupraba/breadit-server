package rawmessageparser

import (
	"encoding/json"
	"time"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/internal/db_types"
)

func ParseJsonComment(data json.RawMessage) (database.Comment, error) {
	var parsed map[string]interface{}

	err := json.Unmarshal(data, &parsed)

	if err != nil {
		return database.Comment{}, err
	}

	dbComment, err := transformJsonComment(parsed)

	return dbComment, err
}

func ParseJsonComments(data json.RawMessage) ([]database.Comment, error) {
	var parsed []map[string]interface{}

	err := json.Unmarshal(data, &parsed)

	if err != nil {
		return []database.Comment{}, err
	}

	comments := []database.Comment{}

	for _, comment := range parsed {
		if comment == nil || int32(comment["id"].(float64)) == 0 {
			continue
		}

		dbComment, err := transformJsonComment(comment)

		if err != nil {
			return []database.Comment{}, err
		}

		comments = append(comments, dbComment)
	}

	return comments, nil
}

func transformJsonComment(commentMap map[string]interface{}) (database.Comment, error) {
	createdAt, err := time.Parse(time.RFC3339, commentMap["created_at"].(string)+"Z")
	updatedAt, err := time.Parse(time.RFC3339, commentMap["updated_at"].(string)+"Z")

	if err != nil {
		return database.Comment{}, err
	}

	replyToId := db_types.NullInt32{}

	if toId, ok := commentMap["reply_to_id"].(float64); ok {
		replyToId.Int32 = int32(toId)
		replyToId.Valid = true
	}

	comment := database.Comment{
		ID:        int32(commentMap["id"].(float64)),
		PostID:    int32(commentMap["post_id"].(float64)),
		AuthorID:  int32(commentMap["author_id"].(float64)),
		Text:      commentMap["text"].(string),
		ReplyToID: replyToId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return comment, nil
}
