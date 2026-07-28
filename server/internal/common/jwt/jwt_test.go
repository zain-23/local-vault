package jwt

import (
	"testing"
	"time"
)

func TestGenerateAndValidate(t *testing.T) {
	svc := NewService("test-secret", 15*time.Minute)

	token, err := svc.GenerateAccessToken("usr_abc", "test@example.com", "")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}

	if claims.Subject != "usr_abc" {
		t.Errorf("subject: want 'usr_abc', got %q", claims.Subject)
	}
	if claims.Email != "test@example.com" {
		t.Errorf("email: want 'test@example.com', got %q", claims.Email)
	}
	if claims.Type != "access" {
		t.Errorf("type: want 'access', got %q", claims.Type)
	}
}

func TestExpiredToken(t *testing.T) {
	svc := NewService("test-secret", 0)
	token, _ := svc.GenerateAccessToken("usr_abc", "test@example.com", "")
	time.Sleep(time.Second) // wait so token is definitely expired

	_, err := svc.ValidateToken(token)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestWrongSecret(t *testing.T) {
	svc1 := NewService("secret-one", 15*time.Minute)
	svc2 := NewService("secret-two", 15*time.Minute)

	token, _ := svc1.GenerateAccessToken("usr_abc", "test@example.com", "")

	// token signed with secret-one should fail with secret-two
	_, err := svc2.ValidateToken(token)
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}

func TestTempToken(t *testing.T) {
	svc := NewService("test-secret", 15*time.Minute)
	token, _ := svc.GenerateTempToken("usr_abc")

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if claims.Type != "2fa_temp" {
		t.Errorf("type: want '2fa_temp', got %q", claims.Type)
	}
}