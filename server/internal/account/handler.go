package account

import (
	"github.com/gofiber/fiber/v2"

	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/middleware"
	"github.com/zain-23/local-vault/server/internal/common/response"
	"github.com/zain-23/local-vault/server/internal/common/validate"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GetMe — GET /api/v1/account/me
func (h *Handler) GetMe(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	res, err := h.svc.Me(c.UserContext(), user.ID)
	if err != nil {
		return err
	}
	return response.Success(c, res, fiber.StatusOK, "profile retrieved")
}

// UpdateProfile — PUT /api/v1/account/me
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c)
	res, err := h.svc.UpdateProfile(c.UserContext(), user.ID, req)
	if err != nil {
		return err
	}
	return response.Success(c, res, fiber.StatusOK, "profile updated")
}

// ChangePassword — PUT /api/v1/account/password
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c)
	if err := h.svc.ChangePassword(c.UserContext(), user.ID, user.SessionID, req); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "password changed")
}

// Enable2FA — POST /api/v1/account/2fa/enable
func (h *Handler) Enable2FA(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	res, err := h.svc.Enable2FA(c.UserContext(), user.ID)
	if err != nil {
		return err
	}
	return response.Success(c, res, fiber.StatusOK, "scan the QR code, then verify")
}

// Verify2FA — POST /api/v1/account/2fa/verify
func (h *Handler) Verify2FA(c *fiber.Ctx) error {
	var req Verify2FARequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c)
	res, err := h.svc.Verify2FA(c.UserContext(), user.ID, req)
	if err != nil {
		return err
	}
	return response.Success(c, res, fiber.StatusOK, "2FA enabled — save your backup codes")
}

// Disable2FA — POST /api/v1/account/2fa/disable
func (h *Handler) Disable2FA(c *fiber.Ctx) error {
	var req Disable2FARequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	user := middleware.GetUser(c)
	if err := h.svc.Disable2FA(c.UserContext(), user.ID, req); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "2FA disabled")
}

// ListSessions — GET /api/v1/account/sessions
func (h *Handler) ListSessions(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	res, err := h.svc.ListSessions(c.UserContext(), user.ID, user.SessionID)
	if err != nil {
		return err
	}
	return response.Success(c, res, fiber.StatusOK, "sessions retrieved")
}

// RevokeSession — DELETE /api/v1/account/sessions/:id
func (h *Handler) RevokeSession(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if err := h.svc.RevokeSession(c.UserContext(), user.ID, c.Params("id")); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "session revoked")
}

// RevokeOtherSessions — DELETE /api/v1/account/sessions
func (h *Handler) RevokeOtherSessions(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if err := h.svc.RevokeOtherSessions(c.UserContext(), user.ID, user.SessionID); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "other sessions revoked")
}
