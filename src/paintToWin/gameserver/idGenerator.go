package main

import (
	"crypto/sha1"
	"fmt"
	"time"
)

func StartIdGenerator() <-chan string {
	idGenChan := make(chan string)

	hasher := sha1.New()
	hasher.Write([]byte(string(time.Now().UnixNano())))

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
