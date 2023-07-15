package Middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/younesabouali/formal-challenges/pet-finder/internal/database"
	"github.com/younesabouali/formal-challenges/pet-finder/utils"
)

var ErrNoAuthHeaderIncluded = errors.New("no authorization header included")

type Middlewares struct {
	DB *database.Queries
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (m *Middlewares) Auth(auth authedHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := GetAPIKey(r.Header)
		if err != nil {
			utils.RespondWithError(w, 403, err.Error())
			return
		}
		user, err := m.DB.GetUserByApiKey(context.Background(), apiKey)
		if err != nil {

			utils.RespondWithError(w, 403, err.Error())
			return
		}
		auth(w, r, user)
		return

	}

}
func (m *Middlewares) IsAdmin(auth authedHandler) func(w http.ResponseWriter, r *http.Request) {
	return m.Auth(
		func(w http.ResponseWriter, r *http.Request, user database.User) {

			if user.Role != "admin" {
				utils.RespondWithError(w, 403, "Unauthorized")
				return
			}
			auth(w, r, user)

		})
}
