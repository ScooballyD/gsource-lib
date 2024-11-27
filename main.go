package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ScooballyD/gsource-lib/internal/database"
	"github.com/ScooballyD/gsource-lib/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("========starting========")
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("unable to open database: %v", err)
	}
	dbQ := database.New(db)

	server.ResetGamebp(dbQ)     //Can disable these if restarting program without wiping the database
	server.ResetDiscountbp(dbQ) //This one especially can take quite awhile

	server.StartServer(dbQ)
}
