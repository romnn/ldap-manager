package grpc

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/romnn/go-service/pkg/auth"
)

// TestTokensAreValid tests if signed JWT tokens can be validated
func TestTokensAreValid(t *testing.T) {
	keyConfig := auth.KeyConfig{
		Generate: true,
	}
	authenticator := auth.Authenticator{
		ExpiresAfter: 1 * time.Hour,
	}

	if err := authenticator.SetupKeys(&keyConfig); err != nil {
		t.Fatalf("failed to setup keys with config %+v: %v", keyConfig, err)
	}

	claims := AuthClaims{
		Username:    "test-user",
		UID:         0,
		IsAdmin:     false,
		DisplayName: "Test User",
	}
	token, err := authenticator.SignJwtClaims(&claims)
	if err != nil {
		t.Fatalf("failed to sign token claims: %+v: %v", claims, err)
	}

	valid, validToken, err := authenticator.Validate(
		token,
		&AuthClaims{},
	)
	t.Log(validToken)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			t.Errorf("malformed token: %v", err)
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			t.Errorf("token is expired: %v", err)
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			t.Errorf("token is not yet valid: %v", err)
		} else {
			t.Errorf("token is bad: %v", err)
		}
	}

	if !valid {
		t.Error("token is not valid")
	}
}
