package game

import (
	"testing"
)

func TestPlayerJoin(t *testing.T) {
	playerOut := make(chan Message, 1)
	player := NewPlayer("testing", true, "testing", "", playerOut)
	game := NewGame("testGame")
	game.addPlayer(&player)

	if game.Players[0] != &player {
		t.Error("Player was not added to the game")
	}
}

func TestSendsBroadcastOnJoin(t *testing.T) {
	player0Out := make(chan Message, 5)
	player1Out := make(chan Message, 5)
	player0 := NewPlayer("testing1", true, "testing1", "", player0Out)
	player1 := NewPlayer("testing2", true, "testing2", "", player1Out)

	game := NewGame("testGame")
	game.addPlayer(&player0)
	game.addPlayer(&player1)

	<-player0Out
	sentMessage := <-player0Out

	joinMessage := NewPlayerJoinMessage("testing2")
	if sentMessage != joinMessage {
		t.Error("Game did not send a correct player join message")
	}
}

func TestLeaveGame(t *testing.T) {
	player0Out := make(chan Message, 5)
	player1Out := make(chan Message, 5)
	player0 := NewPlayer("testing1", true, "testing1", "", player0Out)
	player1 := NewPlayer("testing2", true, "testing2", "", player1Out)

	game := NewGame("testGame")
	game.addPlayer(&player0)

	game.addPlayer(&player1)

	<-player0Out //consume the join message

	game.removePlayer(&player1)

	if len(game.Players) > 1 {
		t.Error("Player was not removed from game")
	}

	<-player0Out
	sentMessage := <-player0Out
	leaveMessage := NewPlayerLeaveMessage("testing2")
	if sentMessage != leaveMessage {
		t.Error("Game did not sent correct player leave message")
	}
}

func TestChat(t *testing.T) {
	player0Out := make(chan Message, 5)
	player1Out := make(chan Message, 5)
	player0 := NewPlayer("testing1", true, "testing1", "", player0Out)
	player1 := NewPlayer("testing2", true, "testing2", "", player1Out)

	game := NewGame("testGame")
	game.addPlayer(&player0)
	game.addPlayer(&player1)

	<-player0Out
	<-player0Out
	<-player1Out

	game.playerChat(&player0, "", "This is test chat")

	chatMessage := NewChatMessage("testing1", "", "This is test chat")

	player0Chat := <-player0Out
	player1Chat := <-player1Out

	if player0Chat != chatMessage {
		t.Error("Invalid player0 chat message sent")
	}

	if player1Chat != chatMessage {
		t.Error("Invalid player1 chat message sent")
	}
}
