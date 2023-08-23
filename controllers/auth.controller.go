package controllers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type signupBody struct {
	Name     string `json:"name" validate:"required,min=2,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type authResponse struct {
	User               models.User `json:"user"`
	AccessToken        string      `json:"accessToken"`
	AccessTokenExpiry  int         `json:"accessTokenExpiry"`
	RefreshToken       string      `json:"refreshToken"`
	RefreshTokenExpiry int         `json:"refreshTokenExpiry"`
}

func (ac *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	var body signupBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := lib.DB.FindUserByEmail(r.Context(), body.Email)

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		fmt.Println("existing user db error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Unable to validate email.")
		return
	}

	if user.ID != 0 {
		utils.RespondWithError(w, http.StatusConflict, "User already exists.")
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

	utils.RespondWithJson(w, http.StatusCreated, utils.Json{"message": "User registered successfully"})
}

func (ac *AuthController) Signin(w http.ResponseWriter, r *http.Request) {
	var body signinBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := lib.DB.FindUserByEmail(r.Context(), body.Email)

	if err != nil || user.ID == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "User does not exist.")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials.")
		return
	}

	accessToken, err := getAccessToken(user, w, r)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := getRefreshToken(user, w, r)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	setLoggedInCookie(w)

	res := authResponse{
		User:               models.DbUserToUser(user),
		AccessToken:        accessToken,
		AccessTokenExpiry:  int(time.Now().Add(constants.AccessTokenTTL).UnixMilli()),
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: int(time.Now().Add(constants.RefreshTokenTTL).UnixMilli()),
	}

	utils.RespondWithJson(w, http.StatusCreated, res)
}

func (ac *AuthController) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	errMessage := "Could not refresh access token"
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, errMessage)
		return
	}

	user, err := utils.GetUserFromToken(w, r, cookie.Value)

	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, "The user belonging to this token no logger exists")
		return
	}

	accessToken, err := getAccessToken(user, w, r)

	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, errMessage)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, utils.Json{"accessToken": accessToken})
}

func (ac *AuthController) GetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	utils.RespondWithJson(w, http.StatusOK, user)
}

func (ac *AuthController) LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, "Could not get refresh token")
		return
	}

	lib.Redis.Del(r.Context(), cookie.Value)

	clearCookies(w)

	utils.RespondWithJson(w, http.StatusOK, struct{}{})
}

func getAccessToken(user database.User, w http.ResponseWriter, r *http.Request) (string, error) {
	accessToken, err := utils.SignJwtToken(strconv.Itoa(int(user.ID)), time.Now().Add(constants.AccessTokenTTL).Unix())

	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(constants.AccessTokenTTL),
		MaxAge:   int(constants.AccessTokenTTL) / int(time.Second),
		HttpOnly: true,
		Secure:   constants.UseSecureCookies,
	})

	return accessToken, nil
}

func getRefreshToken(user database.User, w http.ResponseWriter, r *http.Request) (string, error) {
	refreshToken, err := utils.SignJwtToken(strconv.Itoa(int(user.ID)), time.Now().Add(constants.AccessTokenTTL).Unix())

	if err != nil {
		return "", err
	}

	// ignore the redis error
	err = lib.Redis.Set(r.Context(), refreshToken, user.ID, constants.RefreshTokenTTL).Err()

	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(constants.RefreshTokenTTL),
		MaxAge:   int(constants.RefreshTokenTTL) / int(time.Second),
		HttpOnly: true,
		Secure:   constants.UseSecureCookies,
	})

	return refreshToken, nil
}

func setLoggedInCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Expires:  time.Now().Add(constants.AccessTokenTTL),
		MaxAge:   int(constants.AccessTokenTTL) / int(time.Second),
		HttpOnly: false,
		Secure:   constants.UseSecureCookies,
	})
}

func clearCookies(w http.ResponseWriter) {
	expires := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Now().UTC().Location())

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  expires,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   constants.UseSecureCookies,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  expires,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   constants.UseSecureCookies,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "logged_in",
		Value:    "",
		Expires:  expires,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   constants.UseSecureCookies,
	})
}
