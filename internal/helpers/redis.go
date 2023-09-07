package helpers

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/adhupraba/breadit-server/internal/db_types"
	"github.com/adhupraba/breadit-server/lib"
)

type RedisCachedPost struct {
	ID             int32                   `json:"id" redis:"id"`
	Title          string                  `json:"title" redis:"title"`
	AuthorUsername string                  `json:"authorUsername" redis:"authorUsername"`
	Content        db_types.NullRawMessage `json:"content" redis:"content"`
	CreatedAt      time.Time               `json:"createdAt" redis:"createdAt"`
}

func SavePostInRedis(ctx context.Context, key string, data RedisCachedPost) error {
	err := lib.Redis.HSet(ctx, key, data).Err()

	return err
}

func GetPostFromRedis(ctx context.Context, key string) (RedisCachedPost, error) {
	data, err := lib.Redis.HGetAll(ctx, key).Result()

	if err != nil {
		return RedisCachedPost{}, nil
	}

	id, _ := strconv.Atoi(data["id"])
	content := db_types.NullRawMessage{}
	createdAt, _ := time.Parse(time.RFC3339, data["createdAt"])

	if data["content"] != "" {
		content = db_types.NullRawMessage{
			RawMessage: json.RawMessage(data["content"]),
			Valid:      true,
		}
	}

	post := RedisCachedPost{
		ID:             int32(id),
		Title:          data["title"],
		AuthorUsername: data["authorUsername"],
		Content:        content,
		CreatedAt:      createdAt,
	}

	return post, nil
}
