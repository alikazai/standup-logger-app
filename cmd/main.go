package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("hello world")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Panic().Err(err)
	}
}
