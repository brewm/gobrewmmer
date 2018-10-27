package main

import (
  "os"
  "log"
  "fmt"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"

  "github.com/gin-gonic/gin"

  conn "github.com/brewm/gobrewmmer/pkg/connections"
  temp "github.com/brewm/gobrewmmer/pkg/temperature"
)

func main() {
  initDB()
  defer conn.BrewmmerDB.Close()


  router := gin.Default()

  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  v1 := router.Group("/v1/")
  sessions := v1.Group("/sessions")
  {
   sessions.GET("/",    temp.AllSession)
   sessions.GET("/:id", temp.SingleSession)
   // v1.GET("/", fetchAllTodo)
   // v1.PUT("/:id", updateTodo)
   // v1.DELETE("/:id", deleteTodo)
  }

  router.Run() // listen and serve on 0.0.0.0:8080
}


func initDB() {
  dbPath := getEnv("DB_PATH", "./brewmmer.db")


  _, fileErr := os.Stat(dbPath)
  if os.IsNotExist(fileErr) {
    log.Fatalf("Database file '%s' doesn't exist.", dbPath)
  }

  var dbErr error
  conn.BrewmmerDB, dbErr = sql.Open("sqlite3", dbPath)
  if dbErr != nil {
    log.Fatal(dbErr)
  }
  fmt.Printf("INFO: Database connection to '%s' was successfull!\n", dbPath)
}

func getEnv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
    return value
  }
  return fallback
}