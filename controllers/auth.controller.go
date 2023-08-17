package controllers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/adhupraba/breadit-server/constants"
	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/lib"
	"github.com/adhupraba/breadit-server/models"
	"github.com/adhupraba/breadit-server/utils"
)

type AuthController struct{}

type signinBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signupBody struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	User         models.User `json:"user"`
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
}

func (ac *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	var body signupBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Unable to parse credentials.")
		return
	}

	user, err := lib.DB.FindUserByEmail(r.Context(), body.Email)

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		fmt.Println("existing user db error", err)
		utils.RespondWithError(w, http.StatusNotFound, "Unable to validate email.")
		return
	}

	if user.ID != "" {
		utils.RespondWithError(w, http.StatusBadRequest, "User already exists.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password.")
		return
	}

	randUsername, err := gonanoid.New()

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating username.")
		return
	}

	image := constants.ProfileImages[rand.Intn(len(constants.ProfileImages))]

	user, err = lib.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
		Username: randUsername,
		Image:    sql.NullString{String: image, Valid: true},
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error registering user.")
		return
	}

	accessToken, refreshToken, err := getAccessAndRefreshTokens(user, w, r)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := authResponse{
		User:         models.DbUserToUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	utils.RespondWithJson(w, http.StatusCreated, res)
}

func (ac *AuthController) Signin(w http.ResponseWriter, r *http.Request) {
	var body signinBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Unable to parse credentials.")
		return
	}

	user, err := lib.DB.FindUserByEmail(r.Context(), body.Email)

	if err != nil || user.ID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "User does not exist.")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials.")
		return
	}

	accessToken, refreshToken, err := getAccessAndRefreshTokens(user, w, r)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := authResponse{
		User:         models.DbUserToUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	utils.RespondWithJson(w, http.StatusCreated, res)
}

func getAccessAndRefreshTokens(user database.User, w http.ResponseWriter, r *http.Request) (accessToken string, refreshToken string, err error) {
	accessToken, err = utils.SignJwtToken(user.ID, time.Now().Add(constants.AccessTokenTTL).Unix())

	if err != nil {
		return "", "", err
	}

	refreshToken, err = utils.SignJwtToken(user.ID, time.Now().Add(constants.AccessTokenTTL).Unix())

	if err != nil {
		return "", "", err
	}

	// ignore the redis error
	err = lib.Redis.Set(r.Context(), refreshToken, user.ID, constants.RefreshTokenTTL).Err()

	if err != nil {
		return "", "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(constants.AccessTokenTTL),
		MaxAge:   int(constants.AccessTokenTTL) / int(time.Second),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(constants.RefreshTokenTTL),
		MaxAge:   int(constants.RefreshTokenTTL) / int(time.Second),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	})

	return accessToken, refreshToken, nil
}
