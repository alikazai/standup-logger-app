package main

import (
	"log"

	"github.com/alikazai/standup-logger-app/db"
	"github.com/alikazai/standup-logger-app/utils"
)

func main() {
	conn := db.NewDB(utils.GetDatabaseURL())

	if err := db.RunSeedFiles(conn, "db/seed"); err != nil {
		log.Fatal(err)
	}
}
