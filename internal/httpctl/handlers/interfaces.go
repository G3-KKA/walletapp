package handlers

import "github.com/google/uuid"

//go:generate mockery --filename=mock_wallet_usecase.go --name=WalletUsecase --dir=. --structname=MockWalletUsecase --outpkg=mock_handlers
type WalletUsecase interface {
	UpdateWalletBalance(walletID uuid.UUID, amount int) error
	GetWalletBalance(walletID uuid.UUID) (int, error)
}
