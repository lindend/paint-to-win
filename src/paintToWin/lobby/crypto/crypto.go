package crypto

import (
	"bytes"
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

const SaltSize = 16
const NumIterations = 50000
const PasswordKeySize = 32

const SessionIdSize = 32

func IsValidPassword(password string, hashedPassword []byte, salt []byte) bool {
	newPasswordHash := HashPassword(password, salt)
	return bytes.Compare(hashedPassword, newPasswordHash) == 0
}

func GenerateSessionId() string {
	sessionBytes := make([]byte, SessionIdSize)
	rand.Read(sessionBytes)
	return hex.EncodeToString(sessionBytes)
}

func GenerateSalt() []byte {
	salt := make([]byte, SaltSize)
	rand.Read(salt)
	return salt
}

func HashPassword(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, NumIterations, PasswordKeySize, sha256.New)
}
