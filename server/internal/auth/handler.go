package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/response"
	"github.com/zain-23/local-vault/server/internal/common/validate"
)

// Handler holds auth HTTP handlers
type Handler struct {
	service		*Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Signup handles POST/api/v1/auth/signup
func (h *Handler) Signup(ctx *fiber.Ctx) error {
	var req SignupRequest
	// BodyParser reads JSON body into struct - returns error if JSON is malformed
	if err := ctx.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}

	// Validate struct tags - returns one readable message if a field is missing/invalid
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}

	// c.UserContext() carries request-scoped data — timeouts, cancellation, etc.
	// Message-only endpoint — service returns the text, data stays null
	msg, err := h.service.Signup(ctx.UserContext(), req)
	if err != nil {
		return err
	}
	return response.Success(ctx, "", fiber.StatusCreated, msg )
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	// c.IP() and c.Get("User-Agent") = client info for session tracking
	resp, err := h.service.Login(c.UserContext(), req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "login successful")
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	resp, err := h.service.RefreshToken(c.UserContext(), req)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "token refreshed")
}

// Logout handles POST /api/v1/auth/logout
func (h *Handler) Logout(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	if err := h.service.Logout(c.UserContext(), req.RefreshToken); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "logged out")
}

// VerifyEmail handles POST /api/v1/auth/verify-email
func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	var req VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	resp, err := h.service.VerifyEmail(c.UserContext(), req)
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK, "email verified")
}

// ForgotPassword handles POST /api/v1/auth/forgot-password
func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	msg, err := h.service.ForgotPassword(c.UserContext(), req)
	if err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, msg)
}

// ResetPassword handles POST /api/v1/auth/reset-password
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	msg, err := h.service.ResetPassword(c.UserContext(), req)
	if err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, msg)
}

// SendMagicLink handles POST /api/v1/auth/magic-link
func (h *Handler) SendMagicLink(c *fiber.Ctx) error {
	var req MagicLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	msg, err := h.service.SendMagicLink(c.UserContext(), req)
	if err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, msg)
}

// VerifyMagicLink handles POST /api/v1/auth/magic-link/verify
func (h *Handler) VerifyMagicLink(c *fiber.Ctx) error {
	var req MagicLinkVerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.ErrInvalidBody
	}
	if msg := validate.Struct(req); msg != "" {
		return apperror.New(400, msg)
	}
	resp, err := h.service.VerifyMagicLink(c.UserContext(), req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return err
	}
	return response.Success(c, resp, fiber.StatusOK,  "login successful")
}
