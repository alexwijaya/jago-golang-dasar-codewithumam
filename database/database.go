package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

var DB *pgxpool.Pool

func Connect() {
	var err error
	databaseUrl := viper.GetString("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	DB, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}
	log.Println("Database connection established successfully")
}
