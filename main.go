package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"twentyfour.com/server/model"
	"twentyfour.com/server/play"
)

func main() {

	router := gin.Default()
	rooms := make(map[uint]*model.Room)
	router.GET("/join/:roomid", func(c *gin.Context) {
		rid, err := strconv.Atoi(c.Param("roomid"))
		if err != nil {
			log.Println("room id invalid")
			return
		}
		currRoom, ok := rooms[uint(rid)]
		if !ok {
			log.Println("room id not found")
			return
		}
		play.Join(c, currRoom)
		// c.Param("username")
		// play.Join(c, room)
	})
	router.GET("/create", func(c *gin.Context) {

		room := model.NewRoom()
		rooms[room.Id] = room
		log.Println(room.Id)
		go room.Run()
		play.Join(c, room)
	})
	// router.POST("/create")
	// cards := services.GetCardSet()
	// cards = services.Shuffle(cards)
	// c4rd, cards := services.Get4Card(*cards)
	// log.Printf("4 cards : %v , current card set : %v", c4rd, cards)
	// log.Println(model.FalseUnresolve)
	router.Run()
}
