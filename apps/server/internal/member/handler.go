package member

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
	"github.com/zain-23/local-vault/apps/server/internal/common/middleware"
	"github.com/zain-23/local-vault/apps/server/internal/common/response"
	"github.com/zain-23/local-vault/apps/server/internal/common/validate"
)

// Handler holds member HTTP handlers.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/v1/workspaces/:wid/members?page=&limit=&role=&search=
func (h *Handler) List(c *fiber.Ctx) error {
	var q ListMembersQuery
	if err := c.QueryParser(&q); err != nil { // decode the query string into the struct
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(q); msg != "" { // check page/limit bounds + role enum
		return apperror.New(400, msg)
	}
	q.Normalize()                          // fill page/limit defaults (promoted from Params)
	q.Search = strings.TrimSpace(q.Search) // ignore whitespace-only searches

	resp, err := h.service.List(c.UserContext(), c.Params("wid"), q)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "members retrieved")
}

// Invite handles POST /api/v1/workspaces/:wid/members/invite
func (h *Handler) Invite(c *fiber.Ctx) error {
	var req InviteRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c) // the inviter, from the JWT
	resp, err := h.service.Invite(c.UserContext(), c.Params("wid"), user.ID, req)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusCreated, "invite sent")
}

// Join handles POST /api/v1/workspaces/:wid/members/join (authenticated, no RBAC)
func (h *Handler) Join(c *fiber.Ctx) error {
	var req JoinRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c) // the joiner — id + email come from the JWT
	resp, err := h.service.Join(c.UserContext(), c.Params("wid"), user.ID, user.Email, req)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "joined workspace")
}

// ListInvites handles GET /api/v1/workspaces/:wid/members/invites
func (h *Handler) ListInvites(c *fiber.Ctx) error {
	resp, err := h.service.ListInvites(c.UserContext(), c.Params("wid"))
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "invites retrieved")
}

// CancelInvite handles DELETE /api/v1/workspaces/:wid/members/invites/:id
func (h *Handler) CancelInvite(c *fiber.Ctx) error {
	if err := h.service.CancelInvite(c.UserContext(), c.Params("wid"), c.Params("id")); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "invite cancelled")
}

// ChangeRole handles PUT /api/v1/workspaces/:wid/members/:userId/role
func (h *Handler) ChangeRole(c *fiber.Ctx) error {
	var req ChangeRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	resp, err := h.service.ChangeRole(c.UserContext(), c.Params("wid"), c.Params("userId"), req)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "role updated")
}

// RemoveMember handles DELETE /api/v1/workspaces/:wid/members/:userId
func (h *Handler) RemoveMember(c *fiber.Ctx) error {
	if err := h.service.RemoveMember(c.UserContext(), c.Params("wid"), c.Params("userId")); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "member removed")
}
