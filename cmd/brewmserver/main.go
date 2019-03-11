package main

import (
	"context"
	"database/sql"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"

	"github.com/brewm/gobrewmmer/pkg/protocol/grpc"
	"github.com/brewm/gobrewmmer/pkg/service/brewmmer"
)

// RunServer runs gRPC server and HTTP gateway
func main() {
	// Configure logging
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)

	// Get config
	dbPath := getEnv("BREWM_DB_PATH", "./brewmmer.db")
	grpcPort := getEnv("BREWM_GRPC_PORT", "6999")

	// Set up db
	_, fileErr := os.Stat(dbPath)
	if os.IsNotExist(fileErr) {
		log.WithFields(log.Fields{
			"dbPath": dbPath,
		}).Fatal("Database file doesn't exist!")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Database connection failed!")
	}
	defer db.Close()

	log.WithFields(log.Fields{
		"dbPath": dbPath,
	}).Info("Database connection was successfull!")

	// Start GRPC server
	ctx := context.Background()

	recipeAPI := brewmmer.NewRecipeServiceServer(db)
	sessionAPI := brewmmer.NewSessionServiceServer(db)
	temperatureAPI := brewmmer.NewTemperatureServiceServer()

	grpc.RunServer(ctx, recipeAPI, sessionAPI, temperatureAPI, grpcPort)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
