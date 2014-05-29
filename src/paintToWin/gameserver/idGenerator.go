package main

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"time"
)

func StartIdGenerator() <-chan string {
	idGenChan := make(chan string)

	hasher := sha1.New()
	seed := make([]byte, 32)
	rand.Read(seed)
	hasher.Write(seed)

	counter := 0

	go func() {
		for {
			hasher.Write([]byte(string(time.Now().UnixNano()) + string(counter)))
			counter = counter + 1
			hash := hasher.Sum(nil)
			idGenChan <- fmt.Sprintf("%x", hash)
		}
	}()

	return idGenChan
}
