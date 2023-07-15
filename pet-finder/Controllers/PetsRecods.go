package Controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/younesabouali/formal-challenges/pet-finder/Middlewares"
	"github.com/younesabouali/formal-challenges/pet-finder/internal/database"
	"github.com/younesabouali/formal-challenges/pet-finder/utils"
)

type EventsRecorderController struct {
	DB *database.Queries
}

type RecordParams struct {
	Description string
	ImageUrl    string
	PetId       uuid.UUID

	EventTime string
	Latitude  float64
	Longitude float64
	Area      [][2]float64
}

func (c EventsRecorderController) recordEvent(eventName string, body RecordParams, user database.User) (database.PetsRecord, error) {

	defaultParams := utils.GetDefaultParams()
	LostAt, err := time.Parse(time.RFC3339, body.EventTime)
	if err != nil {
		return database.PetsRecord{}, err
	}

	point := fmt.Sprintf("POINT(%f %f)", body.Longitude, body.Latitude)

	var coords []string
	polygonStr := ""
	if len(body.Area) > 0 {

		for _, point := range body.Area {
			coord := fmt.Sprintf("%f %f", point[0], point[1])
			coords = append(coords, coord)
		}

		polygonStr = "POLYGON((" + strings.Join(coords, ", ") + "))"
		fmt.Println(polygonStr)
	}

	createdRecord := database.PetsRecord{}
	if polygonStr != "" {

		createdRecord, err = c.DB.RecordEvent(context.Background(), database.RecordEventParams{
			ID:        defaultParams.Id,
			Createdat: defaultParams.CreatedAt,
			Updatedat: defaultParams.UpdatedAt,

			PetID:  body.PetId,
			UserID: user.ID,

			ImageUrl:    sql.NullString{String: body.ImageUrl, Valid: true},
			EventName:   database.EventNames(eventName),
			Description: sql.NullString{Valid: true, String: body.Description},

			EventTime:     sql.NullTime{Time: LostAt, Valid: true},
			EventLocation: point,
			Area:          polygonStr,
		})
	} else {

		createdRecord, err = c.DB.RecordEvent(context.Background(), database.RecordEventParams{
			ID:        defaultParams.Id,
			Createdat: defaultParams.CreatedAt,
			Updatedat: defaultParams.UpdatedAt,

			PetID:  body.PetId,
			UserID: user.ID,

			ImageUrl:    sql.NullString{String: body.ImageUrl, Valid: true},
			EventName:   database.EventNames(eventName),
			Description: sql.NullString{Valid: true, String: body.Description},

			EventTime:     sql.NullTime{Time: LostAt, Valid: true},
			EventLocation: point,
		})
	}
	if err != nil {
		fmt.Println(err)
		return createdRecord, err
	}
	return createdRecord, nil
}

func parseRecordParams(r *http.Request) (RecordParams, error) {

	body, err := utils.BodyParser(r, RecordParams{})

	if err != nil {
		return RecordParams{}, errors.New("invalid params")
	}

	if err != nil {
		return RecordParams{}, errors.New("invalid params")
	}
	return body, nil

}

func (c EventsRecorderController) FoundPet(w http.ResponseWriter, r *http.Request, user database.User) {

	eventName := "found"
	body, err := parseRecordParams(r)
	if err != nil {

		utils.RespondWithError(w, 400, err.Error())
		return
	}
	if body.ImageUrl == "" {

		utils.RespondWithError(w, 400, "You need to provide an image for finding pet")
		return
	}
	createdRecord, err := c.recordEvent(eventName, body, user)
	if err != nil {
		utils.RespondWithError(w, 400, "Could record sighting")
		return
	}
	utils.RespondWithJSON(w, 200, createdRecord)

}
func (c EventsRecorderController) FoundPetConfirmation(w http.ResponseWriter, r *http.Request, user database.User) {

	eventName := "rewarded"
	body, err := parseRecordParams(r)
	if err != nil {

		utils.RespondWithError(w, 400, err.Error())
		return
	}
	createdRecord, err := c.recordEvent(eventName, body, user)
	if err != nil {
		utils.RespondWithError(w, 400, "Could record sighting")
		return
	}
	_, err = c.DB.SetPetAsFound(context.Background(), database.SetPetAsFoundParams{ID: body.PetId, UserID: user.ID})

	if err != nil {
		utils.RespondWithError(w, 400, "Could record sighting")
		return
	}
	// add some reward
	utils.RespondWithJSON(w, 200, createdRecord)
}
func (c EventsRecorderController) SightedPet(w http.ResponseWriter, r *http.Request, user database.User) {
	eventName := "sighted"
	body, err := parseRecordParams(r)
	if err != nil {

		utils.RespondWithError(w, 400, err.Error())
		return
	}
	createdRecord, err := c.recordEvent(eventName, body, user)
	if err != nil {
		utils.RespondWithError(w, 400, "Could record sighting")
		return
	}
	utils.RespondWithJSON(w, 200, createdRecord)

}
func (c EventsRecorderController) CommentOnMisingPet(w http.ResponseWriter, r *http.Request, user database.User) {

	eventName := "commentary"
	body, err := parseRecordParams(r)
	if err != nil {

		utils.RespondWithError(w, 400, err.Error())
		return
	}
	createdRecord, err := c.recordEvent(eventName, body, user)
	if err != nil {
		utils.RespondWithError(w, 400, "Could record sighting")
		return
	}
	utils.RespondWithJSON(w, 200, createdRecord)
}
func (c EventsRecorderController) AreaChecked(w http.ResponseWriter, r *http.Request, user database.User) {

	eventName := "area_checked"
	body, err := parseRecordParams(r)
	if err != nil {

		utils.RespondWithError(w, 400, err.Error())
		return
	}
	if len(body.Area) < 4 {

		utils.RespondWithError(w, 400, "a valid area checked need to be provided")
		return
	}
	createdRecord, err := c.recordEvent(eventName, body, user)
	if err != nil {
		utils.RespondWithError(w, 400, "Could record sighting")
		return
	}
	utils.RespondWithJSON(w, 200, createdRecord)
}
func (c EventsRecorderController) GetMissingPetsEvent(w http.ResponseWriter, r *http.Request, user database.User) {
	PetId, err := uuid.Parse(utils.UrlParamsParser(r, "PetId"))
	if err != nil {
		utils.RespondWithError(w, 400, "Unable to parse Params")
		return
	}
	events, err := c.DB.GetMissingPetsEvent(context.Background(), PetId)
	utils.RespondWithArrayJSON(w, 200, events)
	return
}

func EventsRecorderRouter(DB *database.Queries, middlewares Middlewares.Middlewares) *chi.Mux {
	eventRecorderController := EventsRecorderController{DB}
	router := InitializeDependencies(DB)
	router.Post("/SightedPet", middlewares.Auth(eventRecorderController.SightedPet))
	router.Post("/AreaChecked", middlewares.Auth(eventRecorderController.AreaChecked))
	router.Post("/CommentOnMisingPet", middlewares.Auth(eventRecorderController.CommentOnMisingPet))
	router.Post("/FoundPet", middlewares.Auth(eventRecorderController.FoundPet))
	router.Post("/FoundPetConfirmation", middlewares.Auth(eventRecorderController.FoundPetConfirmation))
	router.Get("/GetMissingPetsEvent", middlewares.Auth(eventRecorderController.GetMissingPetsEvent))
	return router

}
