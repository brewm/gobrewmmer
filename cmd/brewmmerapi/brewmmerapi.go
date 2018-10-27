package main

import (
  "github.com/gin-gonic/gin"

  temp "github.com/brewm/gobrewmmer/pkg/temperature"
)

func main() {
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
