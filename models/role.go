package models

type MarketTokenRoles string

const (
	MarketTokenRoleInput      MarketTokenRoles = "INPUT"
	MarketTokenRoleOutput     MarketTokenRoles = "OUTPUT"
	MarketTokenRoleUnderlying MarketTokenRoles = "UNDERLYING"
	MarketTokenRoleLP         MarketTokenRoles = "LP"
	MarketTokenRoleCollateral MarketTokenRoles = "COLLATERAL"
	MarketTokenRoleDebt       MarketTokenRoles = "DEBT"
	MarketTokenRoleReward     MarketTokenRoles = "REWARD"
)
