package game

import (
	"strings"
)

type Round struct {
	word           string
	isFinished     bool
	roundScore     map[*Player]int
	correctGuesses []*Player
	drawing        Drawing
}

func NewRound(word string) Round {
	return Round{word, false, make(map[*Player]int), []*Player{}, NewDrawing()}
}

func (round *Round) Finish() {
	round.isFinished = true
}

func (round *Round) Guess(player *Player, guess string) (wasCorrect bool, hint string) {
	if round.isFinished {
		return false, ""
	}

	if guess == round.word {
		round.correctGuesses = append(round.correctGuesses, player)
		return true, ""
	} else {
		guessHint := ""
		if strings.Contains(round.word, guess) && len(guess)*2 > len(round.word) {
			guessHint = "Your guess is close!"
		}
		return false, guessHint
	}
}
