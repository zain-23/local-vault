package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler, oauth *OAuthHandler, authMW fiber.Handler) {
	// app.Group creates a routes prefix - all routes inside get "/api/v1/auth"
	auth := app.Group("api/v1/auth")

	auth.Post("/refresh", h.RefreshToken)
	auth.Post("/logout", h.Logout)

	// OAuth — :provider is a URL param, accessed via c.Params("provider")
	auth.Get("/oauth/:provider", oauth.RedirectToProvider)
	auth.Get("/oauth/:provider/callback", oauth.HandleCallback)
}
