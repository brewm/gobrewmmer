package main

import (
	"context"
	"database/sql"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"

	"github.com/brewm/gobrewmmer/cmd/brewmserver/global"
	"github.com/brewm/gobrewmmer/pkg/protocol/grpc"
	"github.com/brewm/gobrewmmer/pkg/service/brewmmer"
)

func init() {
	initLogger()
	initDB()
}

// RunServer runs gRPC server and HTTP gateway
func main() {
	defer global.BrewmDB.Close()

	ctx := context.Background()

	recepieAPI := brewmmer.NewRecepieServiceServer()
	sessionAPI := brewmmer.NewSessionServiceServer()

	// listen and serve on localhost:6999
	grpc.RunServer(ctx, recepieAPI, sessionAPI, "6999")
}

func initLogger() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)
}

func initDB() {
	dbPath := getEnv("BREWM_DB_PATH", "./brewmmer.db")

	_, fileErr := os.Stat(dbPath)
	if os.IsNotExist(fileErr) {
		log.WithFields(log.Fields{
			"dbPath": dbPath,
		}).Fatal("Database file doesn't exist!")
	}

	var dbErr error
	global.BrewmDB, dbErr = sql.Open("sqlite3", dbPath)
	if dbErr != nil {
		log.WithFields(log.Fields{
			"err": dbErr,
		}).Fatal("Database connection failed!")
	}

	log.WithFields(log.Fields{
		"dbPath": dbPath,
	}).Info("Database connection was successfull!")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
