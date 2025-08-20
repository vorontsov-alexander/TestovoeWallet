package service

import "github.com/shopspring/decimal"

type WalletBody struct {
	WalletId      string          `json:"wallet_id"`
	OperationType string          `json:"operation_type"`
	Amount        decimal.Decimal `json:"Amount"`
}
