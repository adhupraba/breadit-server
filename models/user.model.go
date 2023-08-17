package models

import (
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/adhupraba/breadit-server/internal/database"
)

type User struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Username  string      `json:"username"`
	Image     null.String `json:"image,omitempty"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

func DbUserToUser(dbUser database.User) User {
	return User{
		ID: dbUser.ID,
		Name: dbUser.Name,
		Email: dbUser.Email,
		Username: dbUser.Username,
		Image: null.NewString(dbUser.Image.String, dbUser.Image.Valid),
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}