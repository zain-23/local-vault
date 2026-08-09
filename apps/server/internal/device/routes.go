package device

import "github.com/gofiber/fiber/v2"

// RegisterRoutes wires device endpoints under /api/v1/device.

func RegisterRoutes(app *fiber.App, h *Handler, authMW fiber.Handler) {
	d := app.Group("/api/v1/device")

	// --- public (CLI, pre-authentication) ---
	d.Post("/authorize", h.Authorize)
	// Registered before "/authorize/:userCode" so "poll" is never parsed as a user code.
	d.Post("/authorize/poll", h.Poll)
	d.Post("/refresh", h.Refresh)

	// --- authenticated (browser + dashboard) ---
	d.Get("/authorize/:userCode", authMW, h.ApprovalDetails)
	d.Put("/authorize/:userCode", authMW, h.Decide)
	d.Get("/", authMW, h.List)
	d.Delete("/:id", authMW, h.Revoke)
}
