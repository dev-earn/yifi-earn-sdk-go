package models

type ChainType string

type Chain struct {
	ChainID     ChainType `json:"chain_id"`
	ImageURL    string    `json:"image_url"`
	Ecosystem   string    `json:"ecosystem"`
	ExplorerUrl string    `json:"explorer_url"`
}
