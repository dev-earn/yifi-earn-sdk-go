package models

type PlatformType string

type Platform struct {
	PlatformID   PlatformType `json:"platform_id"`
	PlatformName string       `json:"platform_name"`
	ImageURL     string       `json:"image"`
	ChainIDs     []string     `json:"chain_ids"`
	CategoryIDs  []string     `json:"category_ids"`
}
