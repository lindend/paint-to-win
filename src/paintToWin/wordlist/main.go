package wordlist

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"paintToWin/storage"
)

func main() {
	var dbConnectionString string

	flag.StringVar(&dbConnectionString, "db", "", "connection string for the database")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Initializing db")
	database, err := storage.InitializeDatabase(dbConnectionString)
	if err != nil {
		log.Fatal("Unable to initialize db ", err)
		return
	}
}
