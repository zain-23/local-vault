package auth

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/middleware"
	"github.com/zain-23/local-vault/server/internal/common/response"
	"github.com/zain-23/local-vault/server/internal/common/validate"
	"github.com/zain-23/local-vault/server/internal/config"
)

// Handler holds auth HTTP handlers
type Handler struct {
	service		*Service
	cfg			config.Config
}

func NewHandler(service *Service, cfg config.Config) *Handler {
	return &Handler{
		service: service,
		cfg: cfg,
	}
}

// ---------------------------- Cookies Helper ----------------
// Package-level so both Handler (password login) and OAuthHandler reuse them.
// setAccessCookie writes the short-lived access token as an HttpOnly cookie.
func setAccessCookie(c *fiber.Ctx, cfg config.Config, token string) {
	c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",                                 // sent to every route
			HTTPOnly: true,
			Secure:   cfg.Env == "production",
			SameSite: "Lax",
			Expires:  time.Now().Add(cfg.JWTAccessExpiry),
	})
}

// setRefreshCookie writes the long-lived refresh token, scoped to the auth routes only.
func setRefreshCookie(c *fiber.Ctx, cfg config.Config, token string) {
	c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    token,
			Path:     "/api/v1/auth",                      // narrow sco
			HTTPOnly: true,
			Secure:   cfg.Env == "production",
			SameSite: "Lax",
			Expires:  time.Now().Add(cfg.JWTRefreshExpiry),
	})
}

func setAuthCookies(c *fiber.Ctx, cfg config.Config, lr *LoginResponse) {
	setAccessCookie(c, cfg, lr.AccessToken)
	setRefreshCookie(c, cfg, lr.RefreshToken)
}

func (h *Handler) clearAuthCookies(c *fiber.Ctx) {
	expired := time.Now().Add(-time.Hour)
	secure := h.cfg.Env == "production"
	c.Cookie(&fiber.Cookie{
		Name: "access_token", 
		Value: "", Path: "/", 
		SameSite: "Lax", 
		Expires: expired,
	})
	c.Cookie(&fiber.Cookie{
		Name: "refresh_token", 
		Value: "", Path: "/api/v1/auth", 
		HTTPOnly: true, 
		Secure: secure, 
		SameSite: "Lax", 
		Expires: expired,
	})
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
	result, err := h.service.Login(c.UserContext(), req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return err
	}

	if result.Requires2FA {
		return response.Success(c, Login2FARequiredResponse{
			Requires2FA: result.Requires2FA,
			TempToken: result.TempToken,
		}, fiber.StatusOK, "2Fa required")
	}

	setAuthCookies(c, h.cfg, result.Tokens)
	return response.Success(c, result.Tokens.User, fiber.StatusOK, "login successful")
}

// Me handles GET /api/v1/auth/me â returns the caller's own account.
func (h *Handler) Me(c *fiber.Ctx) error {
	authUser := middleware.GetUser(c)
	
	user, err := h.service.Me(c.UserContext(), authUser.ID)
	if err != nil {
		return err
	}

	return response.Success(c, user, fiber.StatusOK, "user retrieved")
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
        return apperror.New(fiber.StatusUnauthorized, "refresh token is missing")
    }
	resp, err := h.service.RefreshToken(c.UserContext(), RefreshRequest{RefreshToken:refreshToken})
	if err != nil {
		return err
	}
	setAccessCookie(c, h.cfg, resp.AccessToken)
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

// VerifyEmail handles POST /api/v1/auth/verify-email?token=xxx
func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	// Token comes from the query string, not the body — it's a value from the email link
	token := c.Query("token")
	fmt.Println(token)
	if token == "" {
		return apperror.New(fiber.StatusBadRequest, "verification token is missing")
	}
	if err := h.service.VerifyEmail(c.UserContext(), token); err != nil {
		return err
	}
	return response.Success(c, nil, fiber.StatusOK, "email verified")
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