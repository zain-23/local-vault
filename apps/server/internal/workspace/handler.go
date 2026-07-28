package workspace

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
	"github.com/zain-23/local-vault/apps/server/internal/common/middleware"
	"github.com/zain-23/local-vault/apps/server/internal/common/response"
	"github.com/zain-23/local-vault/apps/server/internal/common/validate"
)

// Handler holds workspace HTTP handlers
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /api/v1/workspaces
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c) // caller id from the JWT
	resp, err := h.service.Create(c.UserContext(), user.ID, req)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusCreated, "workspace created")
}

// List handles GET /api/v1/workspaces
func (h *Handler) List(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	resp, err := h.service.List(c.UserContext(), user.ID)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "workspaces retrieved")
}

// Get handles GET /api/v1/workspaces/:id
func (h *Handler) Get(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	resp, err := h.service.Get(c.UserContext(), user.ID, c.Params("id"))
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "workspace retrieved")
}

// Update handles PUT /api/v1/workspaces/:id (RBAC owner/admin enforced in routes)
func (h *Handler) Update(c *fiber.Ctx) error {
	var req UpdateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	resp, err := h.service.Update(c.UserContext(), c.Params("id"), req)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "workspace updated")
}

// Delete handles DELETE /api/v1/workspaces/:id (RBAC owner enforced in routes)
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.UserContext(), c.Params("id")); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "workspace deleted")
}
