// Package security provides webhook signature validation strategies.
package security

// SignatureValidator validates a webhook payload against a signature header value.
type SignatureValidator interface {
	Validate(payload []byte, signature string) error
}
