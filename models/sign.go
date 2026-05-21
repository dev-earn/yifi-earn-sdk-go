package models

// ValidateRequest is the exact payload used for partner HMAC signing.
type ValidateRequest struct {
	APIKey   string
	Nonce    int64
	Method   string
	Path     string
	BodyHash string
}
