package main

import (
	"github.com/alikazai/standup-logger-app/db"
	"github.com/alikazai/standup-logger-app/utils"
)

func main() {
	db.RunMigrations(utils.GetDatabaseURL(), "db/migrations")
}
