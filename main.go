package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"twentyfour.com/server/model"
	"twentyfour.com/server/play"
)

func main() {
	play.ConnectDb()
	router := gin.Default()
	rooms := make(map[uint64]*model.Room)
	router.GET("/join/:roomid", func(c *gin.Context) {
		rid, err := strconv.Atoi(c.Param("roomid"))
		if err != nil {
			log.Println("room id invalid")
			return
		}
		room, ok := rooms[uint64(rid)]
		if !ok {
			log.Println("room id not found")
			return
		}
		play.Join(c, room)

		// c.Param("username")
		// play.Join(c, room)
	})
	router.GET("/create", func(c *gin.Context) {

		room := model.NewRoom()
		go room.Run()
		rooms[room.Id] = room
		log.Println(room.Id)
		play.Join(c, room)
		// play.Join(c, room)
	})
	// router.POST("/create")
	// cards := services.GetCardSet()
	// cards = services.Shuffle(cards)
	// c4rd, cards := services.Get4Card(*cards)
	// log.Printf("4 cards : %v , current card set : %v", c4rd, cards)
	// log.Println(model.FalseUnresolve)
	// log.Print(services.EvaluateFormula("(8 * 3) * 1 * 2"))
	router.GET("/leaderboard", play.GetLeaderboard)
	router.POST("/leaderboard", play.PostLeaderboard)
	router.Run()
}
