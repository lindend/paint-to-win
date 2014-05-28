package user

import (
	"errors"
	"paintToWin/lobby/crypto"
	"paintToWin/storage"
	"strings"
)

var EmailAlreadyInUseError = errors.New("Email is already in use")
var UsernameAlreadyTakenError = errors.New("Username already taken")
var InvalidUsernameOrPassword = errors.New("Invalid username or password")

func CreateAccount(store *storage.Storage, name string, email string, password string) error {
	salt := crypto.GenerateSalt()
	passwordHash := crypto.HashPassword(password, salt)

	player := storage.Player{
		UserName:     name,
		Email:        email,
		PasswordHash: passwordHash,
		Salt:         salt,
	}
	err := store.Save(&player)
	return err
}

func Login(store *storage.Storage, email string, password string) (storage.Session, error) {
	player := &storage.Player{}
	err := store.FirstWhere(&storage.Player{Email: strings.ToLower(email)}, player)

	if err != nil {
		return storage.Session{}, InvalidUsernameOrPassword
	}

	if !crypto.IsValidPassword(password, player.PasswordHash, player.Salt) {
		return storage.Session{}, InvalidUsernameOrPassword
	}

	session := storage.Session{
		Id:     crypto.GenerateSessionId(),
		Player: player,
	}

	store.SaveInCache("session:"+session.Id, session, 3600)

	return session, nil
}
