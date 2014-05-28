package game

import (
	"testing"
)

func TestFinishRound(t *testing.T) {
	round := NewRound("test")

	if round.isFinished {
		t.Error("Round is finished before Finish being called")
	}

	round.Finish()
	if !round.isFinished {
		t.Error("Expected round.isFinished true")
	}
}

func TestCorrectGuess(t *testing.T) {
	player := NewPlayer("testing", true, "testing", "testing", nil)
	round := NewRound("test")

	wasCorrect, hint := round.Guess(&player, "test")

	if !wasCorrect {
		t.Error("Expected guess to be correct, but was not")
	}

	if hint != "" {
		t.Error("Correct guess should not return a hint")
	}
}

func TestIncorrectGuess(t *testing.T) {
	player := NewPlayer("testing", true, "testing", "testing", nil)
	round := NewRound("test")

	wasCorrect, hint := round.Guess(&player, "banana")

	if wasCorrect {
		t.Error("Incorrect guess returned as correct")
	}

	if hint != "" {
		t.Error("Guess not eligible for hint, but received one anyway")
	}
}

func TestHintableGuess(t *testing.T) {
	player := NewPlayer("testing", true, "testing", "testing", nil)
	round := NewRound("testWord")

	wasCorrect, hint := round.Guess(&player, "testW")
	if wasCorrect {
		t.Error("Guess marked as correct despite being incorrect (just close)")
	}
	if hint == "" {
		t.Error("Guess eligible for hint, but none received")
	}
}

func TestTooShortHintableGuess(t *testing.T) {
	player := NewPlayer("testing", true, "testing", "testing", nil)
	round := NewRound("thisisalongtestword")

	_, hint := round.Guess(&player, "this")
	if hint != "" {
		t.Error("A hint was given but guess is too short to be eligible")
	}
}

func TestCorrectGuessIsRegistered(t *testing.T) {
	player := NewPlayer("testing", true, "testing", "testing", nil)
	round := NewRound("test")

	round.Guess(&player, "test")

	if len(round.correctGuesses) != 1 {
		t.Error("No correct guesses registered despite there being one")
	}

	if round.correctGuesses[0] != &player {
		t.Error("The first correct guess was not assigned to the correct player")
	}
}

func TestIncorrectGuessIsNotRegistered(t *testing.T) {
	player := NewPlayer("testing", true, "testing", "testing", nil)
	round := NewRound("test")

	round.Guess(&player, "fish")

	if len(round.correctGuesses) > 0 {
		t.Error("Some correct guess was registered but no correct guess was made")
	}
}

func TestCorrectGuessOrder(t *testing.T) {
	player1 := NewPlayer("testing1", true, "testing1", "testing", nil)
	player2 := NewPlayer("testing2", true, "testing2", "testing", nil)

	round := NewRound("test")

	round.Guess(&player1, "test")
	round.Guess(&player2, "test")

	if len(round.correctGuesses) != 2 {
		t.Error("Incorrect round.correctGuesses, expected 2")
	}

	if round.correctGuesses[0] != &player1 {
		t.Error("The first player was not the first correct result")
	}

	if round.correctGuesses[1] != &player2 {
		t.Error("The second player was not the second correct result")
	}
}

func TestFinishedRoundAlwaysIncorrect(t *testing.T) {
	player := NewPlayer("testing", true, "testing", "testing", nil)

	round := NewRound("test")

	round.Finish()

	guessResult, _ := round.Guess(&player, "test")

	if guessResult {
		t.Error("Guess was marked correct despite round being finished first")
	}
}
