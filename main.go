package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ScooballyD/gsource-lib/internal/database"
	"github.com/ScooballyD/gsource-lib/scrapers"
	"github.com/ScooballyD/gsource-lib/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("========starting========") //test
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("unable to open database: %v", err)
	}
	dbQ := database.New(db)

	//work in prog
	_, err = scrapers.EpicDeals()
	if err != nil {
		fmt.Println("Unable to collect deals: ", err)
	} //

	server.StartServer(dbQ)
}
