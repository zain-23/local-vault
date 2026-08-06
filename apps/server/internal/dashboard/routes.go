package dashboard

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/apps/server/internal/common/middleware"
)

// RegisterRoutes wires dashboard endpoints under /api/v1/workspaces/:wid/dashboard.
func RegisterRoutes(app *fiber.App, h *Handler, ws middleware.MembershipChecker, authMW fiber.Handler) {
	d := app.Group("/api/v1/workspaces/:wid/dashboard", authMW)
	member := middleware.RequireRole(ws, "wid", RoleOwner, RoleAdmin, RoleMember)

	d.Get("/summary", member, h.Summary)
	d.Get("/activity", member, h.Activity)
}
