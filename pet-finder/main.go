package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	DbManager "github.com/younesabouali/formal-challenges/pet-finder/internal"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	db := DbManager.Manager()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port is not found")
	}
	AppRouter(port, db)
}
