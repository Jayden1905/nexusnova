package api

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/jayden1905/nexusnova/cmd/pkg/database"
	"github.com/jayden1905/nexusnova/config"
	"github.com/jayden1905/nexusnova/service/user"
)

type apiConfig struct {
	addr string
	db   *database.Queries
}

func NewAPIServer(addr string, db *sql.DB) *apiConfig {
	return &apiConfig{
		addr: addr,
		db:   database.New(db),
	}
}

func (s *apiConfig) Run() error {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.Envs.PublicHost,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: true,
	}))

	// Define the api group
	api := app.Group("/api/v1")

	// Define the user store and handler
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)

	// Register the user routes
	userHandler.RegisterRoutes(api)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Error check
	app.Get("/error", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error Checked"})
	})

	log.Println("API Server is running on: ", s.addr)
	return app.Listen(s.addr)
}
