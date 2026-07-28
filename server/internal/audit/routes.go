package audit

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/server/internal/common/middleware"
)

// RegisterRoutes wires audit endpoints under /api/v1/workspaces/:wid/audit.
// Read-only — audit has no create/update/delete over HTTP; events arrive only
// through the Recorder interface from other domains.
func RegisterRoutes(app *fiber.App, h *Handler, ws middleware.MembershipChecker, authMW fiber.Handler) {
	a := app.Group("/api/v1/workspaces/:wid/audit", authMW)

	// List: any member can view the workspace log.
	a.Get("/", middleware.RequireRole(ws, "wid", RoleOwner, RoleAdmin, RoleMember), h.List)
	// Export: owner or admin only.
	a.Get("/export", middleware.RequireRole(ws, "wid", RoleOwner, RoleAdmin), h.Export)
}
