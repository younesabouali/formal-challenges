package Controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/younesabouali/formal-challenges/pet-finder/Middlewares"
	"github.com/younesabouali/formal-challenges/pet-finder/internal/database"
	"github.com/younesabouali/formal-challenges/pet-finder/utils"
)

type StatsController struct {
	DB *database.Queries
}

type Stats struct {
	Percentage    float64
	LostPetCount  int64
	ActiveSpotter int64
}

func (c StatsController) GetStats(w http.ResponseWriter, r *http.Request) {

	percentageCh := make(chan float64)
	lostPetCountCh := make(chan int64)
	activeSpotterCh := make(chan int64)
	errCh := make(chan error)

	// Run the functions in separate goroutines
	go func() {
		percentage, err := c.DB.GetFoundPetPercentage(context.Background())
		if err != nil {
			errCh <- err
			return
		}
		percentageCh <- percentage
	}()
	go func() {
		lostPetCount, err := c.DB.GetLostPetsCount(context.Background())
		if err != nil {
			errCh <- err
			return
		}
		lostPetCountCh <- lostPetCount
	}()
	go func() {
		activeSpotter, err := c.DB.GetActiveSpotter(context.Background())
		if err != nil {
			errCh <- err
			return
		}
		activeSpotterCh <- activeSpotter
	}()

	// Wait for all the goroutines to complete and receive the results or error
	var percentage float64
	var lostPetCount, activeSpotter int64

	for i := 0; i < 3; i++ {
		select {
		case percentage = <-percentageCh:
		case lostPetCount = <-lostPetCountCh:
		case activeSpotter = <-activeSpotterCh:
		case err := <-errCh:
			fmt.Println(err)
			utils.RespondWithError(w, 400, "Cound't compute stats")
			return
		}
	}

	// Create a stats struct with the results
	stats := Stats{
		Percentage:    percentage,
		LostPetCount:  lostPetCount,
		ActiveSpotter: activeSpotter,
	}
	utils.RespondWithJSON(w, 200, stats)
	return
}

func BackofficeRouter(DB *database.Queries, midMiddlewares Middlewares.Middlewares, store *sessions.CookieStore) *chi.Mux {
	statsController := StatsController{DB}
	router := InitializeDependencies(DB)
	router.Use(
		Middlewares.SSOAuthMiddleware(store),
	)
	router.Get("/", statsController.GetStats)
	return router

}
