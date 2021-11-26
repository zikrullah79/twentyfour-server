package model

const (
	GameLog = iota
	PlayerLog
)

const (
	EndGame = iota
	KeyCorrect
	KeyUncorrect
	FalseUnresolve
	PointIncrease
	PlayerJoining
	PlayerKnowSolution
	PostQuestion
	GetLastPlayer
	PostPointedPlayer
	WaitingPointedPlayer
	PlayerLeave
	PausingGame
)
const (
	NewQuestion = iota
	GameStopped
	GameNotStarted
	LastPlayerFound
	WaitingPlayerClaim
	GamePaused
)
const (
	PlayerPointing = iota
	ClaimUnresolve
	StartGame
	ClaimSolution
	AnswerTheQuestion
)

const (
	PlayerPointed = iota
	LastPlayer
	KnowTheSolution
	Thinking
	Unresolve
	WaitingCard
	WaitingResult
)

type UserRequest struct {
	Type            int
	PlayerLogData   *PlayerLogData
	GamePlayLogData *GameLogData
}
type GameLogData struct {
	Question []*int
}
type PlayerLogData struct {
	Id  uint
	Key string
}

type GameResponseNewQuestion struct {
	Type     int
	Question *[]int
}

type GameResponse struct {
	Type int
	Id   uint
}
type GameResponseEndGame struct {
	Type   int
	Result []*FinalResult
}
type FinalResult struct {
	Id    uint
	Point float64
}
