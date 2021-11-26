package model

import (
	"encoding/json"
	"log"
	"math/rand"
	"sort"

	"twentyfour.com/server/services"
)

type Room struct {
	Id      uint64
	Status  int
	Players map[uint]*Player
	// Broadcast  chan []byte
	Broadcast  chan *UserRequest
	Register   chan *Player
	Unregister chan *Player
}

func NewRoom() *Room {
	return &Room{
		Id:      rand.Uint64(),
		Status:  GameNotStarted,
		Players: make(map[uint]*Player),
		// Broadcast:  make(chan []byte),
		Broadcast:  make(chan *UserRequest),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
	}
}

func (r *Room) Run() {
	var card *[]int
	card = services.GetCardSet()
	services.Shuffle(card)
	var pointedPlayer *Player
	var lastPlayer *Player
	var roomStatus int
	for {
		select {
		case player := <-r.Register:
			if len(r.Players) < 8 {
				r.Players[player.Id] = player

				res, err := json.Marshal(&GameResponse{PlayerJoining, player.Id})
				if err != nil {
					continue
				}
				broadcastMessage(res, r)
			}

		case player := <-r.Unregister:
			var logs []interface{}
			if _, ok := r.Players[player.Id]; ok {
				delete(r.Players, player.Id)
				close(player.Send)
				logs = append(logs, &GameResponse{PlayerLeave, player.Id})
				if len(r.Players) < 2 {
					logs = append(logs, &GameResponse{PausingGame, player.Id})
				}
				res, err := json.Marshal(logs)
				if err != nil {
					continue
				}
				r.Status = GamePaused
				broadcastMessage(res, r)
			}

		case logplay := <-r.Broadcast:
			log.Printf("%v logplay", logplay)
			var logs []interface{}
			switch logplay.Type {
			case StartGame:
				if len(r.Players) < 2 {
					continue
				}
				r.Status = NewQuestion
				roomStatus = WaitingPlayerClaim
				setStateToAllUser(Thinking, r)
				rand4Card, currCard := services.Get4Card(*card)
				card = currCard
				logs = append(logs, &GameResponseNewQuestion{PostQuestion, rand4Card})
			case PlayerPointing:
				if logplay.PlayerLogData.Id == 0 || r.Status != LastPlayerFound {
					continue
				}
				if _, ok := r.Players[logplay.PlayerLogData.Id]; !ok {
					continue
				}
				r.Players[logplay.PlayerLogData.Id].State = PlayerPointed
				pointedPlayer = r.Players[logplay.PlayerLogData.Id]
				roomStatus = WaitingPointedPlayer
				logs = append(logs, &GameResponse{PostPointedPlayer, pointedPlayer.Id})
			case ClaimSolution:
				if logplay.PlayerLogData.Id == 0 {
					continue
				}

				if _, count := checkLastPlayer(r); count > 1 {
					r.Players[logplay.PlayerLogData.Id].State = KnowTheSolution
					logs = append(logs, &GameResponse{PlayerKnowSolution, logplay.PlayerLogData.Id})
				}

				if last, count := checkLastPlayer(r); count == 1 {
					r.Players[last].State = LastPlayer
					lastPlayer = r.Players[last]
					// r.Status = LastPlayerFound
					roomStatus = LastPlayerFound
					logs = append(logs, &GameResponse{GetLastPlayer, last})
				}
			case AnswerTheQuestion:
				if (logplay.PlayerLogData.Id == 0 && logplay.PlayerLogData.Key == "") || r.Status != WaitingPointedPlayer {
					continue
				}

				if r.Players[logplay.PlayerLogData.Id].State != PlayerPointed {
					continue
				}

				if services.EvaluateFormula(logplay.PlayerLogData.Key) != nil {
					pointedPlayer.Point -= 4
					logs = append(logs, &GameResponse{PointIncrease, pointedPlayer.Id})
				} else {
					lastPlayer.Point -= 4
					logs = append(logs, &GameResponse{PointIncrease, lastPlayer.Id})
				}
				if len(*card) > 4 {
					r.Status = NewQuestion
					roomStatus = WaitingPlayerClaim
					setStateToAllUser(Thinking, r)
					rand4Card, currCard := services.Get4Card(*card)
					card = currCard
					logs = append(logs, &GameResponseNewQuestion{PostQuestion, rand4Card})
				} else {
					r.Status = GameStopped
					roomStatus = GameStopped
					setStateToAllUser(WaitingResult, r)
					logs = append(logs, &GameResponseEndGame{EndGame, SortByPoint(r)})
				}

			case ClaimUnresolve:
			default:
				continue
			}
			if len(logs) < 1 {
				continue
			}
			res, err := json.Marshal(logs)
			if err != nil {
				continue
			}
			r.Status = roomStatus
			broadcastMessage(res, r)
		}
	}
}

func checkLastPlayer(r *Room) (uint, int) {
	var last []uint
	for _, v := range r.Players {
		if v.State != KnowTheSolution {
			last = append(last, v.Id)
		}
	}
	if len(last) == 1 {
		return last[0], len(last)
	}
	return 0, len(last)
}
func broadcastMessage(msg []byte, r *Room) {
	for _, player := range r.Players {
		select {
		case player.Send <- msg:
		default:
			// log.Println("uhuy")
			close(player.Send)
			delete(r.Players, player.Id)
		}

	}
}

func setStateToAllUser(state int, room *Room) {
	for _, v := range room.Players {
		v.State = state
	}
}

type Players []*FinalResult

func (p Players) Len() int           { return len(p) }
func (p Players) Less(i, j int) bool { return p[i].Point < p[j].Point }
func (p Players) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortByPoint(r *Room) []*FinalResult {
	var players Players
	for _, v := range r.Players {
		players = append(players, &FinalResult{Id: v.Id, Point: v.Point})
	}
	sort.Sort(players)
	return players
}

// func getLastPlayer(r *Room) *Player {
// 	for _, v := range r.Players {
// 		if v.State == LastPlayer {
// 			return v
// 		}
// 	}
// 	return nil
// }
