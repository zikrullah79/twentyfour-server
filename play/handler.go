package play

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{}

func Join(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			// log.Fatal(err)
			break
		}
		log.Printf("Received : %s", msg)
		if err := conn.WriteMessage(msgType, msg); err != nil {
			// log.Fatal(err)
			break
		}

	}
}
