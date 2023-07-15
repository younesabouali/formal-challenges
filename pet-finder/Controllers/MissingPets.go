package Controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/younesabouali/formal-challenges/pet-finder/Middlewares"
	"github.com/younesabouali/formal-challenges/pet-finder/internal/database"
	"github.com/younesabouali/formal-challenges/pet-finder/utils"
)

type MissingPetsController struct {
	DB *database.Queries
}

func (c MissingPetsController) getNearbyMissingPets(w http.ResponseWriter, r *http.Request) {

	type URLParams struct {
		Limit     int32
		Offset    int32
		Latitude  float64
		Longitude float64
		Distance  float64
	}

	body, err := utils.BodyParser(r, URLParams{})
	if err != nil {
		utils.RespondWithError(w, 400, "Unable to parse params")
		return
	}
	// point := "ST_SetSRID(ST_MakePoint(42.123450, -71.987650), 4326)"
	// point := fmt.Sprintf("ST_SetSRID(ST_MakePoint(%f, %f), 4326)", body.Longitude, body.Latitude)
	// fmt.Println(point)
	missingPets, err := c.DB.GetMissingPets(context.Background(), database.GetMissingPetsParams{StMakepoint: body.Latitude, StMakepoint_2: body.Longitude, LostIn: body.Distance, Limit: body.Limit, Offset: body.Offset})
	if err != nil {
		println(err.Error())
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

		Longitude float64
		Latitude  float64
		Lost_at   string
	}
	result, err := utils.BodyParser(r, createMissingPetParams{})
	if err != nil {
		utils.RespondWithError(w, 400, "couldn't parse missing Pet Params")
		return
	}
	defaultParams := utils.GetDefaultParams()

	LostAt, err := time.Parse(time.RFC3339, result.Lost_at)
	if err != nil {

		utils.RespondWithError(w, 400, "Invalid date")
		return
	}

	point := fmt.Sprintf("POINT(%f %f)", result.Longitude, result.Latitude)
	createdMissingPet, err := c.DB.CreateMissingPet(context.Background(), database.CreateMissingPetParams{
		ID:        defaultParams.Id,
		Createdat: defaultParams.CreatedAt,
		Updatedat: defaultParams.UpdatedAt,

		Status:      "missing",
		PetName:     result.PetName,
		Description: sql.NullString{Valid: true, String: result.Description},

		LostIn: point,

		LostAt:   LostAt,
		ImageUrl: sql.NullString{Valid: true, String: result.ImageUrl},

		UserID: user.ID,
	})
	if err != nil {

		fmt.Println(err)
		utils.RespondWithError(w, 400, "Couldn't save the missing pet")
		return
	}
	utils.RespondWithJSON(w, 200, createdMissingPet)

}
func MissingPetsRouter(DB *database.Queries, middlewares Middlewares.Middlewares) *chi.Mux {
	missingPetController := MissingPetsController{DB}
	router := InitializeDependencies(DB)
	router.Post("/newMissingPet", middlewares.Auth(missingPetController.createMessingPetHandler))

	router.Post("/nearbyMissingPets", missingPetController.getNearbyMissingPets)
	// router.Get("/", userController.Login)
	return router

}
