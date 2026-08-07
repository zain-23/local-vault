package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
	"github.com/zain-23/local-vault/apps/server/internal/common/response"
	"github.com/zain-23/local-vault/apps/server/internal/config"
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
// Package-level so both Handler (refresh/logout) and OAuthHandler reuse them.
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
	return response.Success(c, "", fiber.StatusOK, "token refreshed")
}

// Logout handles POST /api/v1/auth/logout — reads the refresh cookie (HttpOnly),
// deletes the session when present, and always clears auth cookies.
func (h *Handler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken != "" {
		if err := h.service.Logout(c.UserContext(), refreshToken); err != nil {
			return err
		}
	}
	h.clearAuthCookies(c)
	return response.Success(c, nil, fiber.StatusOK, "logged out")
}
