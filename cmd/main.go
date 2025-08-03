package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alikazai/standup-logger-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("hello world")
	app := fiber.New()

	// ======================================
	go func() {
		if err := app.Listen(utils.GetHTTPListenAddress()); err != nil {
			log.Panic().Err(err).Msg("Fiber server errror")
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down fiber...")

	if err := app.Shutdown(); err != nil {
		log.Panic().Err(err).Msg("Failed to shut down fiber gracefully")
	}

	log.Info().Msg("Fiber shut down cleanly")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Panic().Err(err)
	}
}
