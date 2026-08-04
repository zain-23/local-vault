package dashboard

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
	"github.com/zain-23/local-vault/apps/server/internal/common/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Summary — GET /api/v1/workspaces/:wid/dashboard/summary
func (h *Handler) Summary(c *fiber.Ctx) error {
	resp, err := h.svc.Summary(c.UserContext(), c.Params("wid"))
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "dashboard summary retrieved")
}

// Activity — GET /api/v1/workspaces/:wid/dashboard/activity
func (h *Handler) Activity(c *fiber.Ctx) error {
	var q ActivityQuery
	if err := c.QueryParser(&q); err != nil {
		return apperror.ErrInvalidBody
	}
	resp, err := h.svc.Activity(c.UserContext(), c.Params("wid"), q)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "dashboard activity retrieved")
}
