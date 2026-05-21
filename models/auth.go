package models

// PartnerAuth stores only the values needed to generate partner headers.
// Nonce and signature are generated per request.
type PartnerAuth struct {
	PartnerID string
	APIKey    string
	Secret    string
}
