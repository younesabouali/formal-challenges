package Controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/younesabouali/formal-challenges/pet-finder/internal/database"
	"github.com/younesabouali/formal-challenges/pet-finder/utils"
)

type MissingPetsController struct {
	DB *database.Queries
}

func (c MissingPetsController) getMissingPets(w http.ResponseWriter, r *http.Request) {

	type URLParams struct {
		Limit  int
		Offset int
	}
	dbParams := database.GetMissingPetsParams{}
	e, err := utils.ParseInt32(utils.UrlParamsParser(r, "Limit"))
	if err != nil {
		utils.RespondWithError(w, 400, "Couldn't Search Missing pets")
		return
	}
	dbParams.Limit = e
	e, err = utils.ParseInt32(utils.UrlParamsParser(r, "Offset"))
	if err != nil {
		utils.RespondWithError(w, 400, "Unable to parse Params")
		return
	}
	dbParams.Offset = e
	missingPets, err := c.DB.GetMissingPets(context.Background(), dbParams)
	if err != nil {
		fmt.Println(err.Error(), dbParams)
		utils.RespondWithError(w, 404, "Unable to perform search")
		return
	}
	utils.RespondWithArrayJSON(w, 200, missingPets)
}
func (c MissingPetsController) createMessingPetHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	type createMissingPetParams struct {
		PetName     string
		Description string
		ImageUrl    string

		Lost_in [2]float64
		Lost_at string
	}
	result, err := utils.BodyParser(r, createMissingPetParams{})
	if err != nil {
		utils.RespondWithError(w, 400, "couldn't parse missing Pet Params")
		return
	}
	defaultParams := utils.GetDefaultParams()

	LostAt, err := time.Parse(time.RFC3339, result.Lost_at)
	if err != nil {

		utils.RespondWithError(w, 400, "couldn't parse missing Pet Params")
		return
	}

	createdMissingPet, err := c.DB.CreateMissingPet(context.Background(), database.CreateMissingPetParams{
		ID:        defaultParams.Id,
		Createdat: defaultParams.CreatedAt,
		Updatedat: defaultParams.UpdatedAt,

		Status:      "missing",
		PetName:     result.PetName,
		Description: sql.NullString{Valid: true, String: result.Description},

		LostIn:   result.Lost_in,
		LostAt:   LostAt,
		ImageUrl: sql.NullString{Valid: true, String: result.ImageUrl},

		UserID: user.ID,
	})
	utils.RespondWithJSON(w, 200, createdMissingPet)

}
func MissingPetsRouter(DB *database.Queries) *chi.Mux {
	missingPetController := MissingPetsController{DB}
	middlewares, router := InitializeDependencies(DB)
	router.Post("/newMissingPet", middlewares.Auth(missingPetController.createMessingPetHandler))

	router.Get("/", missingPetController.getMissingPets)
	// router.Get("/", userController.Login)
	return router

}
