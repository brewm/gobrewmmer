package main

import (
  "github.com/gin-gonic/gin"

  "database/sql"
  _ "github.com/mattn/go-sqlite3"

  temp "github.com/brewm/gobrewmmer/pkg/temperature"
)

var db *sql.DB

func main() {
  initDB()
  defer db.Close()


  r := gin.Default()

  r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  r.GET("/sense", func(c *gin.Context) {
    m := temp.Sense()
    c.JSON(200, m)
  })

  r.Run() // listen and serve on 0.0.0.0:8080
}


func initDB() {
  db, err := sql.Open("sqlite3", "./brewmmer.db")
  if err != nil {
    log.Fatal(err)
  }
  log.Info("Succesfully connected to the database.")
}