package main

import (
	"github.com/gin-gonic/gin"
	"twentyfour.com/server/play"
)

func main() {
	play.ConnectDb()
	router := gin.Default()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	router.GET("/join", play.Join)
	router.GET("/leaderboard", play.GetLeaderboard)
	router.POST("/leaderboard", play.PostLeaderboard)
	router.POST("/create")
	router.Run()
}
