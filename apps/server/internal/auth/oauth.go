package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/zain-23/local-vault/apps/server/internal/config"
)

// OAuthHandler manages Google/Github login flows
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
			"google": {
				ClientID: cfg.GoogleClientID,
				ClientSecret: cfg.GoogleClientSecret,
				RedirectURL: cfg.GoogleRedirectURL,
				Scopes: []string{"openid", "email", "profile"},
				Endpoint: google.Endpoint,
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

// RedirectToProvier - GET /api/v1/auth/oauth/:provider - send browser to Google/Githib login
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
	case "google": 
		return fetchGoogleUser(client)
	}
	
	return nil, fmt.Errorf("unsupported provider: %s", provider)
}

func fetchGoogleUser(client *http.Client) (*oauthUserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // always close response body to prevent memory leaks

	var data struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	json.NewDecoder(resp.Body).Decode(&data)
	return &oauthUserInfo{id: data.ID, email: data.Email, name: data.Name, avatar: data.Picture}, nil
}
