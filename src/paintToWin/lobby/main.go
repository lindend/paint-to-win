package main

import (
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"paintToWin/lobby/api"
	"paintToWin/storage"
	"time"
)

func main() {
	store, err := storage.NewStorage("user=p2wuser password=devpassword host=xena port=5432 dbname=paint2win sslmode=disable",
		"10.10.0.98:6379")
	if err != nil {
		log.Fatal("Unable to initialize storage ", err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter()

	api.RegisterUserApi(router, store)
	api.RegisterGameApi(router, store)

	log.Fatal(http.ListenAndServe(":8083", router))
}
