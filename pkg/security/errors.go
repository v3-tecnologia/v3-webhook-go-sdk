package security

import "errors"

var (
	// ErrMissingSignature is returned when the signature header/value is empty.
	ErrMissingSignature = errors.New("missing webhook signature")
	// ErrInvalidSignature is returned when the signature does not match the payload.
	ErrInvalidSignature = errors.New("invalid webhook signature")
)
