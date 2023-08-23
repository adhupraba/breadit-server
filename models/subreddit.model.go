package models

import (
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/adhupraba/breadit-server/internal/database"
)

type Subreddit struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatorID null.Int  `json:"creatorId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func DbSubredditToSubreddit(dbSubreddit database.Subreddit) Subreddit {
	return Subreddit{
		ID:        dbSubreddit.ID,
		Name:      dbSubreddit.Name,
		CreatorID: null.NewInt(int64(dbSubreddit.CreatorID.Int32), dbSubreddit.CreatorID.Valid),
		CreatedAt: dbSubreddit.CreatedAt,
		UpdatedAt: dbSubreddit.UpdatedAt,
	}
}
