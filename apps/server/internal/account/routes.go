package account

import "github.com/gofiber/fiber/v2"

// RegisterRoutes wires account endpoints under /api/v1/account (all require a JWT).
func RegisterRoutes(app *fiber.App, h *Handler, authMW fiber.Handler) {
	a := app.Group("/api/v1/account", authMW)

	a.Get("/me", h.GetMe)
	a.Put("/me", h.UpdateProfile)

	a.Get("/sessions", h.ListSessions)
	// DELETE "/sessions" (all others) is registered distinctly from "/sessions/:id".
	a.Delete("/sessions", h.RevokeOtherSessions)
	a.Delete("/sessions/:id", h.RevokeSession)
}
