package security_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/security"
)

func TestHMACSHA256Validator(t *testing.T) {
	secret := "test-secret"
	payload := []byte(`{"id":"1"}`)

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	validator, err := security.NewHMACSHA256Validator(secret)
	if err != nil {
		t.Fatalf("new validator: %v", err)
	}

	if err := validator.Validate(payload, signature); err != nil {
		t.Fatalf("expected valid signature, got %v", err)
	}

	if err := validator.Validate(payload, "DEADBEEF"); !errors.Is(err, security.ErrInvalidSignature) {
		t.Fatalf("expected invalid signature, got %v", err)
	}

	if err := validator.Validate(payload, ""); !errors.Is(err, security.ErrMissingSignature) {
		t.Fatalf("expected missing signature, got %v", err)
	}
}

func TestNewHMACSHA256Validator_EmptySecret(t *testing.T) {
	if _, err := security.NewHMACSHA256Validator("  "); err == nil {
		t.Fatal("expected error for empty secret")
	}
}
