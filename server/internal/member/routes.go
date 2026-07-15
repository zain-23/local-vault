package member

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/server/internal/common/middleware"
)

// RegisterRoutes wires member endpoints under /api/v1/workspaces/:wid/members.
// authMW (JWT) protects every route; RequireRole reads the caller's role from the
// ":wid" param and gates each route. `join` has NO RequireRole — the joiner isn't a member yet.
func RegisterRoutes(app *fiber.App, h *Handler, store *Store, authMW fiber.Handler) {
	m := app.Group("/api/v1/workspaces/:wid/members", authMW)

	// Accept an invite: authenticated only. RequireRole would 404 a non-member.
	m.Post("/join", h.Join)

	// Read the roster: any member (owner, admin, or member).
	m.Get("/", middleware.RequireRole(store, "wid", RoleOwner, RoleAdmin, RoleMember), h.List)

	// Invites: owner or admin.
	m.Post("/invite", middleware.RequireRole(store, "wid", RoleOwner, RoleAdmin), h.Invite)
	m.Get("/invites", middleware.RequireRole(store, "wid", RoleOwner, RoleAdmin), h.ListInvites)
	m.Delete("/invites/:id", middleware.RequireRole(store, "wid", RoleOwner, RoleAdmin), h.CancelInvite)

	// Change a member's role: owner only.
	m.Put("/:userId/role", middleware.RequireRole(store, "wid", RoleOwner), h.ChangeRole)

	// Remove a member: owner or admin.
	m.Delete("/:userId", middleware.RequireRole(store, "wid", RoleOwner, RoleAdmin), h.RemoveMember)
}
