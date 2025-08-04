package main

import (
	"embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alikazai/standup-logger-app/utils"
	"github.com/gofiber/fiber/v2"
	htmlTemplate "github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var viewsFS embed.FS

func main() {
	fmt.Println("hello world")
	engine := htmlTemplate.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Static("/static", "public")

	app.Get("/", func(c *fiber.Ctx) error {
		log.Info().Msg("Homepage")
		return c.Render("index", fiber.Map{
			"Title": "Welcome to Standup logger App",
		})
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		log.Info().Msg("login")
		return c.Render("login", fiber.Map{
			"Title": "Login",
		})
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		log.Info().Msg("register")
		return c.Render("register", fiber.Map{
			"Title": "Register",
		})
	})
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
