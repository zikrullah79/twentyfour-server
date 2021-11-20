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
			if u.PlayerType != RoomMaster || u.Room.Status != GameNotStarted {
				continue
			}

			u.Room.Broadcast <- &UserRequest{StartGame, nil, nil}

		case PlayerPointing:
			if p.PlayerLogData == nil || u.State != LastPlayer {
				continue
			}
			u.Room.Broadcast <- &UserRequest{PlayerPointing, &PlayerLogData{p.PlayerLogData.Id, ""}, nil}
		case ClaimSolution:
			if u.Room.Status != WaitingPlayerClaim || u.State != Thinking {
				continue
			}
			u.Room.Broadcast <- &UserRequest{ClaimSolution, &PlayerLogData{u.Id, ""}, nil}
		case ClaimUnresolve:
			u.State = Unresolve
		case AnswerTheQuestion:
			if u.State != PlayerPointed {
				continue
			}

			u.Room.Broadcast <- &UserRequest{AnswerTheQuestion, &PlayerLogData{u.Id, p.PlayerLogData.Key}, nil}
		}
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
			// que := len(u.Send)
			// log.Println(que)
			// for i := 0; i < que; i++ {
			// 	w.Write([]byte("/n"))
			// 	w.Write(<-u.Send)
			// }
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
