package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// HMACSHA256Validator validates webhook signatures using HMAC-SHA256 hex digests.
type HMACSHA256Validator struct {
	secret []byte
}

// NewHMACSHA256Validator creates a validator for the given UTF-8 secret.
func NewHMACSHA256Validator(secret string) (*HMACSHA256Validator, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, fmt.Errorf("hmac secret cannot be empty")
	}
	return &HMACSHA256Validator{secret: []byte(secret)}, nil
}

// Validate compares signature with the hex-encoded HMAC-SHA256 of payload.
func (v *HMACSHA256Validator) Validate(payload []byte, signature string) error {
	if strings.TrimSpace(signature) == "" {
		return ErrMissingSignature
	}

	mac := hmac.New(sha256.New, v.secret)
	_, _ = mac.Write(payload)
	computed := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(strings.ToLower(computed)), []byte(strings.ToLower(signature))) {
		return ErrInvalidSignature
	}
	return nil
}
