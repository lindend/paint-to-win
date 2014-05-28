package game

import (
	"errors"
	"reflect"
)

var NoSuchHandlerError = errors.New("No such message handler registered")
var HandlerAlreadyRegistered = errors.New("A message handler for this message type is already registered")
var HandlerIsNotFuncError = errors.New("handlerFunction is not a function")
var InvalidHandlerFuncSignatureError = errors.New("Invalid signature of handler function, expected func (*Player, *MessageDataType)")

var playerType = reflect.TypeOf(&Player{})

type Receiver interface{}
type HandlerFunction interface{}

type MessageHandler struct {
	Handlers map[reflect.Type]reflect.Value
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		make(map[reflect.Type]reflect.Value),
	}
}

func (m MessageHandler) Handle(message InMessage) error {
	dataType := reflect.TypeOf(message.Message.Data)
	if handler, exists := m.Handlers[dataType]; !exists {
		return NoSuchHandlerError
	} else {
		dataValue := reflect.ValueOf(message.Message.Data)
		playerValue := reflect.ValueOf(message.Source)
		handler.Call([]reflect.Value{playerValue, dataValue})
	}
	return nil
}

func (m *MessageHandler) Add(handlerFunction HandlerFunction) error {
	handlerType := reflect.TypeOf(handlerFunction)

	handlerKind := handlerType.Kind()
	if handlerKind != reflect.Func {
		return HandlerIsNotFuncError
	}

	numArguments := handlerType.NumIn()
	if numArguments != 2 {
		return InvalidHandlerFuncSignatureError
	}

	if handlerType.In(0) != playerType {
		return InvalidHandlerFuncSignatureError
	}

	messageType := handlerType.In(1)
	if _, exists := m.Handlers[messageType]; exists {
		return HandlerAlreadyRegistered
	}

	m.Handlers[messageType] = reflect.ValueOf(handlerFunction)

	return nil
}
