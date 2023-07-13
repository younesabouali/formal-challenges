package DbManager

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/younesabouali/formal-challenges/pet-finder/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}

func Manager() *database.Queries {
	fmt.Println("New Initialization")

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("dbURL is not found")
	}
	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	return ApiConfig{DB: database.New(conn)}.DB
}
