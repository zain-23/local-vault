package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {
	// app.Group creates a routes prefix - all routes inside get "/api/v1/auth"
	auth := app.Group("api/v1/auth")

	auth.Post("/signup", h.Signup)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)
	auth.Post("/logout", h.Logout)
	auth.Post("/verify-email", h.VerifyEmail)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Post("/reset-password", h.ResetPassword)
	auth.Post("/magic-link", h.SendMagicLink)
	auth.Post("/magic-link/verify", h.VerifyMagicLink)
}
