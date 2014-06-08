package game

import (
	"encoding/json"
	"errors"
)

type Message struct {
	Type string
	Data interface{}
}

type InMessage struct {
	Message Message
	Source  *Player
}

type MessagePlayer struct {
	Id   string
	Name string
}

type WelcomeMessage struct {
	Player  MessagePlayer
	Message string
}

type GameStateMessage struct {
	Players       []MessagePlayer
	State         string
	StateData     interface{}
	TimeRemaining int
}

type TurnToPaintMessage struct {
	PaintingPlayerId string
}

type TurnToChooseWordMessage struct {
	ChoosingPlayerId string
}

type StrokesMessage struct {
	Strokes []Stroke
}

type ChatMessage struct {
	To      string
	From    string
	Message string
}

type GuessMessage struct {
	Guess string
}

type PlayerJoinMessage struct {
	Player MessagePlayer
}

type PlayerLeaveMessage struct {
	PlayerId string
}

//type NewRoundMessage struct {
//	DrawingPlayerId string
//}
//
//type RoundWordMessage struct {
//	Word string
//}

type PlayerScore struct {
	Score    int
	PlayerId string
}

type EndRoundMessage struct {
	Scores      []PlayerScore
	TotalScores []PlayerScore
	CorrectWord string
}

type HintMessage struct {
	Hint      string
	HintLevel int
}

type ChooseWordMessage struct {
	Word string
}

type CorrectGuessMessage struct {
	PlayerId string
}

type CloseGuessMessage struct {
	Guess string
	Hint  string
}

type WrongGuessMessage struct {
	PlayerId string
	Guess    string
}

const MsgId_PlayerJoin = "PlayerJoin"
const MsgId_PlayerLeave = "PlayerLeave"
const MsgId_Chat = "Chat"
const MsgId_GameState = "GameState"

//const MsgId_NewRound = "NewRound"
const MsgId_Strokes = "Strokes"

//const MsgId_TurnToPaint = "TurnToPaint"
//const MsgId_TurnToChooseWord = "TurnToChooseWord"
const MsgId_ChooseWord = "ChooseWord"
const MsgId_CorrectGuess = "CorrectGuess"
const MsgId_CloseGuess = "CloseGuess"
const MsgId_WrongGuess = "WrongGuess"
const MsgId_Guess = "Guess"
const MsgId_Welcome = "Welcome"

func NewWelcomeMessage(player MessagePlayer, message string) Message {
	return Message{MsgId_Welcome, WelcomeMessage{player, message}}
}

func NewPlayerJoinMessage(player MessagePlayer) Message {
	return Message{MsgId_PlayerJoin, PlayerJoinMessage{player}}
}

func NewPlayerLeaveMessage(id string) Message {
	return Message{MsgId_PlayerLeave, PlayerLeaveMessage{id}}
}

func NewChatMessage(from string, to string, message string) Message {
	return Message{MsgId_Chat, ChatMessage{to, from, message}}
}

//func NewNewRoundMessage(drawingPlayerId string) Message {
//	return Message{MsgId_NewRound, NewRoundMessage{drawingPlayerId}}
//}

func NewGameStateMessage(state string, stateData interface{}, players []MessagePlayer, timeLeft int) Message {
	return Message{MsgId_GameState, GameStateMessage{players, state, stateData, timeLeft}}
}

func NewStrokesMessage(strokes []Stroke) Message {
	return Message{MsgId_Strokes, StrokesMessage{strokes}}
}

//func NewTurnToPaintMessage(playerId string) Message {
//	return Message{MsgId_TurnToPaint, TurnToPaintMessage{playerId}}
//}

//func NewTurnToChooseWordMessage(playerId string) Message {
//	return Message{MsgId_TurnToChooseWord, TurnToChooseWordMessage{playerId}}
//}

func NewCorrectGuessMessage(playerId string) Message {
	return Message{MsgId_CorrectGuess, CorrectGuessMessage{playerId}}
}

func NewCloseGuessMessage(guess string, hint string) Message {
	return Message{MsgId_CloseGuess, CloseGuessMessage{guess, hint}}
}

func NewWrongGuessMessage(playerId string, guess string) Message {
	return Message{MsgId_WrongGuess, WrongGuessMessage{playerId, guess}}
}

type internalMessage struct {
	Type string
	Data json.RawMessage
}

func getMessageStruct(messageType string) (interface{}, error) {
	var msgData interface{}
	switch messageType {
	case MsgId_Chat:
		msgData = new(ChatMessage)
	case MsgId_GameState:
		msgData = new(GameStateMessage)
	case MsgId_Strokes:
		msgData = new(StrokesMessage)
	case MsgId_ChooseWord:
		msgData = new(ChooseWordMessage)
	case MsgId_Guess:
		msgData = new(GuessMessage)
	default:
		return nil, errors.New("Unable to find decoder for message type " + messageType)
	}
	return msgData, nil
}

func DecodeMessage(message []byte) (Message, error) {
	internalMsg := internalMessage{}
	err := json.Unmarshal(message, &internalMsg)

	if err != nil {
		return Message{}, err
	}

	var msgData interface{}
	msgData, err = getMessageStruct(internalMsg.Type)
	if err != nil {
		return Message{}, err
	}

	err = json.Unmarshal(internalMsg.Data, msgData)
	if err != nil {
		return Message{}, err
	}

	return Message{internalMsg.Type, msgData}, nil
}

func EncodeMessage(message Message) ([]byte, error) {
	return json.Marshal(message)
}
