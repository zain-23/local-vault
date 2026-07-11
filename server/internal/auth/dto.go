package auth

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
}

type SignupResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type LoginResult struct {
	Requires2FA bool				`json:"requires_2fa"`
	TempToken   string				`json:"temp_token"`
	Tokens      *LoginResponse		`json:"tokens"`
}

type Login2FARequiredResponse struct {
	Requires2FA bool				`json:"requires_2fa"`
	TempToken   string				`json:"temp_token"`
}

type Login2FARequest struct {
	TempToken string `json:"temp_token"`
	TOTPCode  string `json:"totp_code"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type MagicLinkRequest struct {
	Email string `json:"email"`
}

type MagicLinkVerifyRequest struct {
	Token string `json:"token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}