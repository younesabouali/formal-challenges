package Controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	Middlewares "github.com/younesabouali/formal-challenges/pet-finder/Middlewares"
	"github.com/younesabouali/formal-challenges/pet-finder/internal/database"
	"github.com/younesabouali/formal-challenges/pet-finder/utils"
)

type UserController struct {
	DB *database.Queries
}
type Controller interface {
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (c UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type userParams struct {
		Name     string
		Email    string
		Password string
	}
	result, err := utils.BodyParser(r, userParams{})
	if err != nil {
		utils.RespondWithError(w, 400, "couldn't parse user")
		return
	}
	// validate the email and password
	hash, err := HashPassword(result.Password)
	if err != nil {
		utils.RespondWithError(w, 400, "couldn't parse user")
		return
	}
	createdUser, err := c.DB.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(),
		Name:      result.Name,
		Updatedat: time.Now(),
		Createdat: time.Now(),
		Email:     result.Email,
		Role:      "user",
		Password:  hash,
	})
	if err != nil {
		utils.RespondWithError(w, 400, "Couldn't create user")
		return
	}
	utils.RespondWithJSON(w, 200, createdUser.ApiKey)
}
func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func (c UserController) Login(w http.ResponseWriter, r *http.Request) {
	type LoginParams struct {
		Email    string
		Password string
	}

	result, err := utils.BodyParser(r, LoginParams{})

	if err != nil {
		utils.RespondWithError(w, 400, "couldn't parse user")
		return
	}
	foundUser, err := c.DB.FindUserByEmail(context.Background(), result.Email)

	if err != nil {
		if err == sql.ErrNoRows {

			utils.RespondWithError(w, 404, "user Not found")
			return
		}
		utils.RespondWithError(w, 400, "couldn't parse user")
		return
	}

	err = verifyPassword(result.Password, foundUser.Password)
	if err != nil {

		utils.RespondWithError(w, 404, "user not found")
		return
	}
	utils.RespondWithJSON(w, 200, foundUser.ApiKey)
}
func InitializeDependencies(DB *database.Queries) (Middlewares.Middlewares, *chi.Mux) {
	middlewares := Middlewares.Middlewares{DB: DB}
	router := chi.NewRouter()
	return middlewares, router
}
func UserRouter(DB *database.Queries) *chi.Mux {
	userController := UserController{DB}
	_, router := InitializeDependencies(DB)
	router.Post("/registerUser", userController.CreateUserHandler)
	router.Get("/login", userController.Login)
	return router

}
