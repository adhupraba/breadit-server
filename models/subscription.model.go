package models

import (
	"time"

	"github.com/adhupraba/breadit-server/internal/database"
)

type Subscription struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"userId"`
	SubredditID int32     `json:"subredditId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func DbSubscriptionToSubscription(dbSubscription database.Subscription) Subscription {
	return Subscription{
		ID:          dbSubscription.ID,
		UserID:      dbSubscription.UserID,
		SubredditID: dbSubscription.SubredditID,
		CreatedAt:   dbSubscription.CreatedAt,
		UpdatedAt:   dbSubscription.UpdatedAt,
	}
}
