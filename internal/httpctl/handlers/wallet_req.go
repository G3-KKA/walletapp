package handlers

import (
	"github.com/google/uuid"
)

// Generalized error on Wallet operations.
type WalletError struct {
	Error string `json:"error"`
}

type UpdateWalletBalanceRequest struct {
	WalletID      uuid.UUID `json:"valletid"      validate:"required"` // Misspell.
	OperationType string    `json:"operationType" validate:"required,oneof=WITHDRAW DEPOSIT"`
	Amount        uint      `json:"amount"        validate:"gte=0"`
}
type UpdateWalletBalanceResponse struct {
	Msg string `json:"msg"`
}
type GetWalletBalanceRequest struct {
}
type GetWalletBalanceResponse struct {
	Balance int `json:"balance"`
}
