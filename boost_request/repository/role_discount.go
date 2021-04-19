package repository

import "github.com/shopspring/decimal"

type RoleDiscount struct {
	ID       int64
	GuildID  string
	RoleID   string
	Discount decimal.Decimal
}
