package api

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/jayden1905/go-nextjs-template/cmd/pkg/database"
	"github.com/jayden1905/go-nextjs-template/config"
	"github.com/jayden1905/go-nextjs-template/service/user"
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
		AllowOrigins:     config.Envs.FrontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: true,
	}))

	// Define the user store and handler
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)

	// Define the api group
	api := app.Group("/api/v1")

	// Register the user routes
	userHandler.RegisterRoutes(api)

	log.Println("API Server is running on: ", s.addr)
	return app.Listen(s.addr)
}
