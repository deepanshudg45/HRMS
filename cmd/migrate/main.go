package main

import (
	"context"
	"log"

	"WITS/config"
	"WITS/package/database"
	"WITS/package/migration"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migration.Run(context.Background(), db, "migration"); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations completed successfully")
}
