package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/zain-23/local-vault/apps/server/internal/config"
)

// OAuthHandler manages the GitHub login flow
type OAuthHandler struct {
	service 	*Service
	config		map[string]*oauth2.Config
	frontend	string
	cfg			config.Config
}

func NewOAuthHandler(service *Service, cfg config.Config) *OAuthHandler {
	return &OAuthHandler{
		service: service,
		config: map[string]*oauth2.Config{
			"github": {
				ClientID: cfg.GithubClientID,
				ClientSecret: cfg.GithubClientSecret,
				RedirectURL: cfg.GithubRedirectURL,
				Scopes: []string{"read:user", "user:email"},
				Endpoint: github.Endpoint,
			},
		},
		frontend: cfg.FrontendURL,
		cfg: cfg,
	}
}

func generateOAuthState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// RedirectToProvier - GET /api/v1/auth/oauth/:provider - send browser to GitHub login
func (h *OAuthHandler) RedirectToProvider(ctx *fiber.Ctx) error {
	provider := ctx.Params("provider")
	oauthCfg, ok := h.config[provider]

	if !ok {
		return  fiber.NewError(400, "unsupported provider")
	}

	// state token prevent CSRF - we store it in cookie, verify on callback
	state := generateOAuthState()
	ctx.Cookie(&fiber.Cookie{
		Name: "oauth_state",
		Value: state,
		MaxAge: 300,  // 5minutes
		HTTPOnly: true, // javascript can't read it
		SameSite: "Lax", // send on top-level navigation only
	})

	// AuthCodeURL builds the full redirect URL with client_id, scopes, state, etc.
	return ctx.Redirect(oauthCfg.AuthCodeURL(state))
}


// HandleCallback - GET /api/v1/auth/oauth/:provider/callback
func (h *OAuthHandler) HandleCallback(ctx *fiber.Ctx) error {
	provider := ctx.Params("provider")
	authCfg, ok := h.config[provider]

	if !ok {
		return fiber.NewError(400, "unsuported provider")
	}

	// verify state matches cookie - prevent CSRF attacks
	if ctx.Cookies("oauth_state") != ctx.Query("state") {
		return ctx.Redirect(h.frontend + "/auth/login?error=invalid+state")
	}

	// Clear state cookie - one-time use
	ctx.Cookie(&fiber.Cookie{Name: "oauth_state", Value: "", MaxAge: -1})

	// Exchange the one-time authorization CODE for an access token.
	token, err := authCfg.Exchange(context.TODO(), ctx.Query("code"))
	if err != nil {
		return ctx.Redirect(h.frontend + "/auth/login?error=oauth+failed")
	}

	// fetch use info from provider's API using the token
	oauthUser, err := fetchOAuthUser(provider, authCfg, token)
	if err != nil {
		return ctx.Redirect(h.frontend + "/auth/login?error=profile+failed")
	}

	// Find or create OAuth user in DB
	user, err := h.service.FindOrCreateOAuthUser(
		ctx.UserContext(), provider, oauthUser.id, oauthUser.email, oauthUser.name, oauthUser.avatar,
	)

	if err != nil {
		return ctx.Redirect(h.frontend + "/auth/login?error=account+failed")
	}

	// create session and token
	loginResp, err := h.service.OAuthLogin(ctx.UserContext(), user, ctx.IP(), ctx.Get("User-Agent"))

	if err != nil {
		return ctx.Redirect(h.frontend + "/auth/login?error=session+failed")
	}
	
	// Set the same HttpOnly session cookies as password login, then land on the app.
	setAuthCookies(ctx, h.cfg, loginResp)
	return ctx.Redirect(h.frontend + "/dashboard")
}

// ----------------- Helper -------------
type oauthUserInfo struct {
	id, email, name, avatar string
}

// FetchOAuthUser calls provider API to get user profile
func fetchOAuthUser(provider string, cfg *oauth2.Config, token *oauth2.Token) (*oauthUserInfo, error) {
	client := cfg.Client(context.TODO(), token)

	switch provider {
	case "github":
		return fetchGithubUser(client)
	}

	return nil, fmt.Errorf("unsupported provider: %s", provider)
}

// fetchGithubUser calls GitHub's /user API for profile info, then falls back to
// /user/emails when the primary email is private (GitHub omits it from /user in that case).
func fetchGithubUser(client *http.Client) (*oauthUserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // always close response body to prevent memory leaks

	var data struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	email := data.Email
	if email == "" {
		email, err = fetchGithubPrimaryEmail(client)
		if err != nil {
			return nil, err
		}
	}

	name := data.Name
	if name == "" {
		name = data.Login // many GitHub users never set a display name
	}

	return &oauthUserInfo{
		id:     strconv.FormatInt(data.ID, 10),
		email:  email,
		name:   name,
		avatar: data.AvatarURL,
	}, nil
}

// fetchGithubPrimaryEmail hits /user/emails (requires the user:email scope) — needed
// when the user has "Keep my email address private" enabled, so /user returns no email.
func fetchGithubPrimaryEmail(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	var firstVerified string
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
		if e.Verified && firstVerified == "" {
			firstVerified = e.Email
		}
	}
	if firstVerified != "" {
		return firstVerified, nil
	}
	return "", fmt.Errorf("no verified email found on GitHub account")
}
