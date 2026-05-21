package models

type CategoryID string

const (
	DexCategory              CategoryID = "dex"
	LendingCategory          CategoryID = "lending"
	LSTCategory              CategoryID = "lst"
	FixedIncomeCategory      CategoryID = "fixed-income"
	StableCoinsCategory      CategoryID = "stablecoins"
	StructuredVaultsCategory CategoryID = "structured-vaults"
)

type Category struct {
	CategoryID  CategoryID `json:"category_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
}
