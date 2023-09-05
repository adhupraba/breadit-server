package helpers

import (
	"time"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/internal/types"
)

type RedisCachedPost struct {
	ID             int32                 `json:"id" redis:"id"`
	Title          string                `json:"title" redis:"title"`
	AuthorUsername string                `json:"authorUsername" redis:"authorUsername"`
	Content        types.NullRawMessage  `json:"content" redis:"content"`
	CurrentVote    database.NullVoteType `json:"currentVote" redis:"currentVote"`
	CreatedAt      time.Time             `json:"createdAt" redis:"createdAt"`
}
