package user

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/jayden1905/nexusnova/config"
	"github.com/jayden1905/nexusnova/service/auth"
	"github.com/jayden1905/nexusnova/service/email"
	"github.com/jayden1905/nexusnova/types"
	"github.com/jayden1905/nexusnova/utils"
)

type Handler struct {
	store  types.UserStore
	mailer email.Mailer
}

func NewHandler(store types.UserStore, mailer email.Mailer) *Handler {
	return &Handler{store: store, mailer: mailer}
}

// RegisterRoutes for Fiber
func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/auth/login", h.handleLogin)
	router.Post("/auth/register", h.handleRegister)
	// router.Get("/user", auth.WithJWTAuth(h.handleGetUser, h.store))
}

// handler for registering a new user
func (h *Handler) handleRegister(c *fiber.Ctx) error {
	// Parse JSON payload
	var payload types.RegisterUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	invalidFields, validationErr := utils.ValidatePayload(payload)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid payload",
			"invalid_fields": invalidFields,
		})
	}

	// Check if the user already exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("User with email %s already exists", payload.Email)})
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	// Create a new user with unverified status
	err = h.store.CreateUser(c.Context(), &types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate a verification token
	token, err := auth.GenerateVerificationToken(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send email asynchronously
	go func() {
		err = h.mailer.SendVerificationEmail(payload.Email, token)
		if err != nil {
			fmt.Printf("Error sending verification email: %v\n", err)
		}
	}()

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"email":   payload.Email,
		"status":  "verification email sent",
	})
}

func (h *Handler) handleLogin(c *fiber.Ctx) error {
	// Parse JSON payload
	var payload types.LoginUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Invalid payload: %s", errors)})
	}

	// Check if the user exists by email
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or password is incorrect"})
	}

	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or password is incorrect"})
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, int(u.UserID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the token with HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,     // Set to true in production (HTTPS)
		SameSite: "Strict", // Prevent CSRF attacks
		Path:     "/",      // Valid for the entire site
		MaxAge:   int(config.Envs.JWTExpirationInSeconds),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token, "expires_in": fmt.Sprintf("%d", config.Envs.JWTExpirationInSeconds)})
}
