package main

import (
  "os"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"

  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"

  global "github.com/brewm/gobrewmmer/pkg/global"
  temp "github.com/brewm/gobrewmmer/pkg/temperature"
)

func init() {
  initLogger()
  initDB()
}

func main() {
  defer global.BrewmDB.Close()

  router := gin.Default()

  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  v1 := router.Group("/v1/")

  v1.GET("/sense", temp.Sense)
  sessions := v1.Group("/sessions")
  {
    sessions.GET("/",    temp.AllSession)
    sessions.GET("/:id", temp.SingleSession)
    sessions.POST("/",   temp.StartSession)
    sessions.PUT("/",    temp.StopSession)
  }

  router.Run() // listen and serve on 0.0.0.0:8080
}

func initLogger() {
  log.SetLevel(log.InfoLevel)
  log.SetOutput(os.Stderr)
}

func initDB() {
  dbPath := getEnv("DB_PATH", "./brewmmer.db")


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