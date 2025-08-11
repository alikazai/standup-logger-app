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
var database *sql.DB

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}
type StandupEntry struct {
	ID        int
	UserID    uuid.UUID
	Name      string
	Yesterday string
	Today     string
	Blockers  string
	Date      time.Time
}

func main() {
	database = db.NewDB(utils.GetDatabaseURL())
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

		loggedIn := sess.Get("logged_in")
		log.Info().Msg("logged_in")
		if loggedIn == true {
			userEmail := sess.Get("user_email").(string)
			log.Info().Msg(userEmail)
			today := time.Now().UTC()
			var selectedDate time.Time
			var err error

			dateParam := c.Query("date")
			if dateParam != "" {
				selectedDate, err = time.Parse("2006-01-02", dateParam)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).SendString("invalid date format")
				}
			} else {
				selectedDate = today
			}
			entries, err := getEntriesByDate(selectedDate)
			if err != nil {
				log.Info().Msg("NO ENTRIES")
				return err
			}
			return c.Render("index", fiber.Map{
				"Title":        "Welcome to Standup logger App",
				"SelectedDate": selectedDate.Format("2006-01-02"),
				"Today":        today.Format("2006-01-02"),
				"Entries":      entries,
			})
		} else {
			return c.Redirect("/login")
		}
	})

	app.Post("/standup", func(c *fiber.Ctx) error {
		log.Info().Msg("standup form")
		sess, _ := store.Get(c)

		loggedIn := sess.Get("logged_in")

		if loggedIn == true {
			type StandupForm struct {
				Yesterday string `form:"yesterday"`
				Today     string `form:"today"`
				Blockers  string `form:"blockers"`
			}
			var form StandupForm

			if err := c.BodyParser(&form); err != nil {
				return err
			}
			userEmail := sess.Get("user_email").(string)
			user, err := getUserByEmail(userEmail)
			if err != nil {
				log.Error().Err(err).Msg("failed to fetch user")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Something went wrong",
				})
			}

			standup, err := getEntryByUserAndToday(user.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Something went wrong",
				})
			}
			if standup != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "You have already submitted for today",
				})
			}

			if err = submitStandupEntry(user.ID, form.Yesterday, form.Today, form.Blockers); err != nil {
				log.Error().Err(err).Msg("failed to insert standup entry")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Could not add standup entry",
				})
			}

			return c.Redirect("/")
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
		if utils.Password_verify(user.Password, []byte(form.Password)) {
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

		_, err = database.Exec(`INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`,
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
	row := database.QueryRow(`SELECT id, name, email, password, created_at FROM users WHERE email = $1`, email)

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

func getEntriesByDate(date time.Time) ([]*StandupEntry, error) {
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	fmt.Println(start)
	fmt.Println(end)

	rows, err := database.Query(`
  SELECT se.id, se.user_id, u.name, se.yesterday, se.today, se.blockers, se.date
  FROM standup_entries se
  JOIN users u ON se.user_id = u.id
  WHERE se.date >= $1 AND se.date < $2
`, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*StandupEntry
	for rows.Next() {
		var entry StandupEntry
		if err := rows.Scan(&entry.ID, &entry.UserID, &entry.Name, &entry.Yesterday, &entry.Today, &entry.Blockers, &entry.Date); err != nil {
			return entries, err
		}
		entries = append(entries, &entry)
	}

	if err := rows.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}

func getEntryByUserAndToday(userID uuid.UUID) (*StandupEntry, error) {
	row := database.QueryRow(`
  SELECT id, user_id, date
  FROM standup_entries
  WHERE user_id = $1
    AND date >= CURRENT_DATE
    AND date <  CURRENT_DATE + INTERVAL '1 day'
`, userID)

	var entry StandupEntry
	if err := row.Scan(&entry.ID, &entry.UserID, &entry.Date); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // user not found
		}
		return nil, err // other DB error
	}

	return &entry, nil
}

func submitStandupEntry(userID uuid.UUID, yesterday, today, blockers string) error {
	date := time.Now().UTC()
	_, err := database.Exec(`INSERT INTO standup_entries (user_id, yesterday, today, blockers, date) VALUES ($1, $2, $3, $4, $5)`,
		userID, yesterday, today, blockers, date)

	return err
}
