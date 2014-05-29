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

type GameStateMessage struct {
	Players       []string
	State         string
	StateData     interface{}
	TimeRemaining int
}

type TurnToPaintMessage struct {
	PaintingPlayer string
}

type TurnToChooseWordMessage struct {
	ChoosingPlayer string
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
	PlayerName string
}

type PlayerLeaveMessage struct {
	PlayerName string
}

type NewRoundMessage struct {
	DrawingPlayer string
}

type RoundWordMessage struct {
	Word string
}

type PlayerScore struct {
	Score      int
	PlayerName string
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
	PlayerName string
}

type CloseGuessMessage struct {
	Guess string
	Hint  string
}

type WrongGuessMessage struct {
	PlayerName string
	Guess      string
}

const MsgId_PlayerJoin = "PlayerJoin"
const MsgId_PlayerLeave = "PlayerLeave"
const MsgId_Chat = "Chat"
const MsgId_GameState = "GameState"
const MsgId_NewRound = "NewRound"
const MsgId_Strokes = "Strokes"
const MsgId_TurnToPaint = "TurnToPaint"
const MsgId_TurnToChooseWord = "TurnToChooseWord"
const MsgId_ChooseWord = "ChooseWord"
const MsgId_CorrectGuess = "CorrectGuess"
const MsgId_CloseGuess = "CloseGuess"
const MsgId_WrongGuess = "WrongGuess"
const MsgId_Guess = "Guess"

func NewPlayerJoinMessage(name string) Message {
	return Message{MsgId_PlayerJoin, PlayerJoinMessage{name}}
}

func NewPlayerLeaveMessage(name string) Message {
	return Message{MsgId_PlayerLeave, PlayerLeaveMessage{name}}
}

func NewChatMessage(from string, to string, message string) Message {
	return Message{MsgId_Chat, ChatMessage{to, from, message}}
}

func NewNewRoundMessage(drawingPlayer string) Message {
	return Message{MsgId_NewRound, NewRoundMessage{drawingPlayer}}
}

func NewGameStateMessage(state string, stateData interface{}, players []string, timeLeft int) Message {
	return Message{MsgId_GameState, GameStateMessage{players, state, stateData, timeLeft}}
}

func NewStrokesMessage(strokes []Stroke) Message {
	return Message{MsgId_Strokes, StrokesMessage{strokes}}
}

func NewTurnToPaintMessage(paintingPlayer string) Message {
	return Message{MsgId_TurnToPaint, TurnToPaintMessage{paintingPlayer}}
}

func NewTurnToChooseWordMessage(choosingPlayer string) Message {
	return Message{MsgId_TurnToChooseWord, TurnToChooseWordMessage{choosingPlayer}}
}

func NewCorrectGuessMessage(playerName string) Message {
	return Message{MsgId_CorrectGuess, CorrectGuessMessage{playerName}}
}

func NewCloseGuessMessage(guess string, hint string) Message {
	return Message{MsgId_CloseGuess, CloseGuessMessage{guess, hint}}
}

func NewWrongGuessMessage(playerName string, guess string) Message {
	return Message{MsgId_WrongGuess, WrongGuessMessage{playerName, guess}}
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
