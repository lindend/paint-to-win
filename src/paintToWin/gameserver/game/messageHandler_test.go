package game

import (
	"fmt"
	"testing"
)

type TestHandler struct {
}

func (t *TestHandler) HandleMessage(pl *Player, sm *StrokesMessage) {
	fmt.Println("yay")
}

func TestMessageHandler(t *testing.T) {
	h := new(TestHandler)
	pl := new(Player)
	sm := new(StrokesMessage)

	mhandler := NewMessageHandler()
	if err := mhandler.Add(h.HandleMessage); err != nil {
		t.Error(err)
	}
	if err := mhandler.Handle(InMessage{Message{MsgId_Strokes, sm}, pl}); err != nil {
		t.Error(err)
	}
}
