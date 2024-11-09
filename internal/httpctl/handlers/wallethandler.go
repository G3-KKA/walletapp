package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type wallethandler struct {
	l     *zerolog.Logger
	v     *validator.Validate
	ucase WalletUsecase
}

// NewWalletHandler is a wallet handler constructor.
func NewWalletHandler(l *zerolog.Logger, ucase WalletUsecase) *wallethandler {
	return &wallethandler{
		l:     l,
		v:     validator.New(), // Singleton enough?
		ucase: ucase,
	}
}

// nolint:godot // Swagger logic cannot contain dots at the end.
//
// @Summary 	Update wallet balance.
// @Description Update wallet if type of operation and amount are permitted.
// @ID			update
// @Tags 		wallet
// @Accept 		json
// @Produce		json
// @Param		request body UpdateWalletBalanceRequest true "update wallet"
// @Success 	200 {object} UpdateWalletBalanceResponse
// @Failure     400 {object} WalletError
// @Failure     404 {object} WalletError
// @Failure     500 {object} WalletError
// @Router 		/api/v1/wallet [POST]
func (handler *wallethandler) Update(gctx *gin.Context) {

	req := UpdateWalletBalanceRequest{} // nolint:exhaustruct // For unmarshall .

	err := gctx.BindJSON(&req)
	if err != nil {
		handler.l.
			Error().
			Any("wallethandler.Update req:", req). // Unsafe if req contains private info.
			Errs("errs", []error{err, ErrGotInvalidJSON}).
			Send()

		JSONE(gctx, err, http.StatusBadRequest, WalletError{
			Error: ErrGotInvalidJSON.Error(),
		})

		return
	}

	err = handler.v.Struct(req)
	if err != nil {
		handler.l.
			Error().
			Any("wallethandler.Update req:", req).
			Errs("errs", []error{err, ErrUnsuccessfulValidation}).
			Send()

		JSONE(gctx, err, http.StatusBadRequest, WalletError{
			Error: ErrUnsuccessfulValidation.Error(),
		})

		return
	}
	amountConverter := func(delta int, optype string) int {
		switch optype {
		case "WITHDRAW":
			return -delta
		case "DEPOSIT":
			return delta
		default:

			return 0
		}
	}

	err = handler.ucase.UpdateWalletBalance(
		req.WalletID,
		amountConverter(int(req.Amount), req.OperationType), // nolint:gosec // Impossible overflow.
	)

	if err != nil {
		handler.l.
			Error().
			Any("wallethandler.Update req:", req).
			Err(err).Err(ErrInernalError).
			Send()
		JSONE(gctx, err, http.StatusInternalServerError, WalletError{
			Error: ErrInernalError.Error(),
		})

		return
	}

	gctx.JSON(http.StatusOK, UpdateWalletBalanceResponse{
		Msg: MessageSuccess,
	})

}

// nolint:godot // Swagger logic cannot contain dots at the end.
//
// @Summary 	Get wallet balance.
// @Description Get wallet balance if that wallet exists.
// @ID 			get
// @Tags		wallet
// @Produce 	json
// @Param 		wallet_uuid	 path string true "identifier"
// @Success     200 {object} GetWalletBalanceResponse
// @Failure     400 {object} WalletError
// @Failure     404 {object} WalletError
// @Failure     500 {object} WalletError
// @Router		/api/v1/wallets/{wallet_uuid} [GET]
func (handler *wallethandler) GetBalance(gctx *gin.Context) {

	walletid, err := uuid.Parse(gctx.Params.ByName("wallet_uuid"))

	if err != nil {
		handler.l.
			Error().
			Str("trace", "wallethandler.GetBalance").
			Err(ErrInvalidUUID).Err(err).Send()

		JSONE(gctx, err, http.StatusBadRequest, WalletError{
			Error: ErrInvalidUUID.Error(),
		})

		return
	}

	balance, err := handler.ucase.GetWalletBalance(walletid)
	if err != nil {
		handler.l.
			Error().
			Str("trace", "wallethandler.Update req").
			Err(err).Err(ErrInernalError).
			Send()

		JSONE(gctx, err, http.StatusInternalServerError, WalletError{
			Error: ErrInernalError.Error(),
		})

		return
	}

	gctx.JSON(http.StatusOK, GetWalletBalanceResponse{
		Balance: balance,
	})

}
