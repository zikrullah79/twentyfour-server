package model

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id    uint
	State int
	Point float64
	Conn  *websocket.Conn
	Room  *Room
	Send  chan []byte
}

func (u *Player) ReadPlayerUpdate() {
	for {
		_, msg, err := u.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var p = &LogHistory{}
		err = json.Unmarshal(msg, p)
		if err != nil {
			log.Println(err)
			continue
		}
		switch p.Type {
		case GameStarted:
			log.Println("Game Started")

		}
		// log.Println(p.Type)
		// log.Println(msg)
		u.Room.Broadcast <- bytes.TrimSpace(msg)
	}
}

func (u *Player) WritePlayerUpdate() {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		u.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-u.Send:
			// log.Printf("%s", msg)
			if !ok {
				u.Conn.NextWriter(websocket.CloseMessage)
				return
			}
			w, err := u.Conn.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Println(err)
				return
			}
			w.Write(msg)
			// var p = &LogHistory{}
			// err = json.Unmarshal(msg, p)
			// if err != nil {
			// 	log.Println(err)
			// 	// break
			// }
			// log.Println(p)
			que := len(u.Send)
			// log.Println(que)
			for i := 0; i < que; i++ {
				w.Write([]byte("/n"))
				w.Write(<-u.Send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := u.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
