package models

import "time"

type TokensRequest struct {
	Search string
	Chains []string
	Tags   []string
	Group  string
	Limit  int64
	Offset *int64
}

type GetTokensResponse struct {
	TokenID string `json:"token_id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Image   string `json:"image,omitempty"`
	GroupID string `json:"token_group"`

	TagIDs   []string `json:"tag_ids"`
	ChainIDs []string `json:"chain_ids"`
}

type FromTo struct {
	From time.Time
	To   time.Time
}
