package play

import (
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"twentyfour.com/server/model"
)

var Upgrader = websocket.Upgrader{}

// func Join(c *gin.Context) {
// 	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	defer conn.Close()

// 	for {
// 		msgType, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			// log.Fatal(err)
// 			break
// 		}
// 		log.Printf("Received : %s", msg)
// 		if err := conn.WriteMessage(msgType, msg); err != nil {
// 			// log.Fatal(err)
// 			break
// 		}

// 	}
// }

func Join(c *gin.Context, r *model.Room) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}

	// defer conn.Close()

	newPlayer := &model.Player{
		Id:    uint(rand.Uint64()),
		State: model.PlayerJoining,
		Point: 0,
		Conn:  conn,
		Room:  r,
		Send:  make(chan []byte, 256),
	}
	// log.Printf("%v room", &r)
	log.Printf("%v room", r.Players)
	// log.Printf("%v joined", newPlayer.Id)
	newPlayer.Room.Register <- newPlayer
	go newPlayer.WritePlayerUpdate()
	go newPlayer.ReadPlayerUpdate()
}
