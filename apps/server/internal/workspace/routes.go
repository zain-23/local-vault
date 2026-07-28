package workspace

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/apps/server/internal/common/middleware"
)

// RegisterRoutes wires workspace endpoints under /api/v1/workspaces.
// authMW protects every route; PUT/DELETE additionally require a role.
func RegisterRoutes(app *fiber.App, h *Handler, store *Store, authMW fiber.Handler) {
	// Passing authMW as a second arg to Group applies it to every route below
	ws := app.Group("/api/v1/workspaces", authMW)

	ws.Post("/", h.Create)
	ws.Get("/", h.List)
	ws.Get("/:id", h.Get)

	// RequireRole reads the workspace id from the ":id" param, then checks the caller's role.
	// Update: owner OR admin. Delete: owner only.
	ws.Put("/:id", middleware.RequireRole(store, "id", RoleOwner, RoleAdmin), h.Update)
	ws.Delete("/:id", middleware.RequireRole(store, "id", RoleOwner), h.Delete)
}
