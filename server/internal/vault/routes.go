package vault

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/server/internal/common/middleware"
)

func RegisterRoutes(app *fiber.App, h *Handler, ws middleware.MembershipChecker, authMW fiber.Handler) {
	v := app.Group("/api/v1/workspaces/:wid/vaults", authMW)

	// Any member (owner/admin/member) can read, create vaults/tokens, and sync.
	member := middleware.RequireRole(ws, "wid", RoleOwner, RoleAdmin, RoleMember)
	v.Post("/", member, h.Create)
	v.Get("/", member, h.List)
	v.Get("/:id", member, h.Get)
	v.Put("/:id/snapshot", member, h.PushSnapshot)
	v.Get("/:id/snapshot", member, h.PullSnapshot)
	v.Post("/:id/tokens", member, h.CreateToken)
	v.Get("/:id/tokens", member, h.ListTokens)
	v.Delete("/:id/tokens/:tid", member, h.RevokeToken)

	// Destructive: owner or admin only.
	admin := middleware.RequireRole(ws, "wid", RoleOwner, RoleAdmin)
	v.Delete("/:id", admin, h.Delete)
	v.Delete("/:id/peers/:did", admin, h.RemovePeer)

	// --- top-level ---
	app.Post("/api/v1/join", h.Join) // public

	msg := app.Group("/api/v1/messages", authMW)
	msg.Post("/", h.SendMessage)
	msg.Get("/:deviceId", h.GetMessages)
}
