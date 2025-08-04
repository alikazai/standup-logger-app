package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alikazai/standup-logger-app/db"
	"github.com/alikazai/standup-logger-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	htmlTemplate "github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var store = session.New()
var databse *sql.DB

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func main() {
	databse = db.NewDB(utils.GetDatabaseURL())
	fmt.Println("hello world")
	engine := htmlTemplate.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	app.Use(recover.New())
	app.Static("/static", "public")

	app.Get("/", func(c *fiber.Ctx) error {
		log.Info().Msg("Homepage")
		sess, _ := store.Get(c)

		logged_in := sess.Get("logged_in")
		log.Info().Msg("logged_in")
		if logged_in == true {
			userEmail := sess.Get("user_email").(string)
			log.Info().Msg(userEmail)
			// user_data := models.GetUserDataByEmail(userEmail)
			return c.Render("index", fiber.Map{
				"Title": "Welcome to Standup logger App",
			})
		} else {
			return c.Redirect("/login")
		}
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		log.Info().Msg("login")
		sess, _ := store.Get(c)
		errorMsg := sess.Get("login_error")
		sess.Delete("login_error")
		sess.Save()

		return c.Render("login", fiber.Map{
			"Title":       "Login",
			"login_error": errorMsg,
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		type SignupForm struct {
			Email    string `form:"email"`
			Password string `form:"password"`
		}
		var form SignupForm

		if err := c.BodyParser(&form); err != nil {
			return err
		}

		// get user from db
		user, err := getUserByEmail(form.Email)
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch user")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
		}
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
		}
		userVerified := false
		if utils.Password_verify(user.Password, []byte(form.Password)) == true {
			log.Error().Err(err).Msg("Login successful")
			userVerified = true
		} else {
			log.Error().Err(err).Msg("Login not successful")
			return c.Status(fiber.StatusUnauthorized).SendString("Login not successful")
		}

		if userVerified {
			sess, _ := store.Get(c)
			sess.Set("user_email", user.Email)
			sess.Set("logged_in", true)

			sess.Save()

			return c.Redirect("/")
		}
		return c.Redirect("/login")
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		log.Info().Msg("register")
		return c.Render("register", fiber.Map{
			"Title": "Register",
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		type RegisterForm struct {
			Name     string `form:"name"`
			Email    string `form:"email"`
			Password string `form:"password"`
		}

		var form RegisterForm

		if err := c.BodyParser(&form); err != nil {
			log.Error().Err(err).Msg("failed to parse form")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form submission")
		}

		log.Info().Str("email", form.Email).Msg("registration submitted")

		// TODO: save to DB
		hashedPassword, err := utils.HashPassword(form.Password)
		if err != nil {
			log.Error().Err(err).Msg("failed to hash password")
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
		}

		_, err = databse.Exec(`INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`,
			form.Name, form.Email, hashedPassword)

		if err != nil {
			log.Error().Err(err).Msg("failed to insert user")
			return c.Status(fiber.StatusInternalServerError).SendString("Could not register user")
		}

		return c.SendString("Registration successful!")
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

func getUserByEmail(email string) (*User, error) {
	row := databse.QueryRow(`SELECT id, name, email, password, created_at FROM users WHERE email = $1`, email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // user not found
		}
		return nil, err // other DB error
	}

	return &user, nil
}
