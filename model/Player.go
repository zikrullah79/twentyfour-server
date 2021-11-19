package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	RoomMember = iota
	RoomMaster
)

type Player struct {
	Id         uint
	State      int
	PlayerType int
	Point      float64
	Conn       *websocket.Conn
	Room       *Room
	Send       chan []byte
}

func (u *Player) ReadPlayerUpdate() {
	for {
		_, msg, err := u.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var p = &UserRequest{}
		err = json.Unmarshal(msg, p)
		if err != nil {
			log.Println(err)
			continue
		}
		switch p.Type {
		case StartGame:
			if u.PlayerType != RoomMaster || u.Room.Status == NewQuestion {
				continue
			}
			// SetStateToAllUser(WaitingCard, u.Room)

			u.Room.Broadcast <- &UserRequest{StartGame, nil, nil}
			// log.Printf("Game Started by %v", u.PlayerType)

		case PlayerPointing:
			if p.PlayerLogData == nil || u.State != LastPlayer {
				continue
			}
			if p, ok := u.Room.Players[p.PlayerLogData.Id]; !ok {
				p.State = PlayerPointed
			}
		case ClaimSolution:
			// if u.State != Thinking {
			// 	continue
			// }
			// u.State = KnowTheSolution
			u.Room.Broadcast <- &UserRequest{ClaimSolution, &PlayerLogData{u.Id, ""}, nil}
		case ClaimUnresolve:
			u.State = Unresolve
		case AnswerTheQuestion:
			if u.State != PlayerPointed {
				continue
			}
		}
		// log.Println(p.Type)
		// log.Println(msg)
		// u.Room.Broadcast <- bytes.TrimSpace(msg)
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

func SetStateToAllUser(state int, room *Room) {
	for _, v := range room.Players {
		v.State = state
	}
}
