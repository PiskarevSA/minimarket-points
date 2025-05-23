package handlers

import (
	"context"

	"github.com/PiskarevSA/minimarket-points/internal/dto"
	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
)

var getBalance500JSONResponse = oapi.GetBalance500JSONResponse{
	Code:    oapi.S1394,
	Message: oapi.InternalError,
}

func (h *Handlers) GetBalance(
	ctx context.Context,
	request oapi.GetBalanceRequestObject,
) (oapi.GetBalanceResponseObject, error) {
	const op = "getBalance"

	token, ok := getJwtFromContext(ctx, op)
	if !ok {
		return getBalance500JSONResponse, nil
	}

	userId, err := getUserIdFromJwt(token, op)
	if err != nil {
		return getBalance500JSONResponse, nil
	}

	balance, err := h.getBalanceUsecase.Do(ctx, userId)
	if err != nil {
		return getBalance500JSONResponse, nil
	}

	getBalance := dto.BalanceToGetBalanceResponse(balance)

	return oapi.GetBalance200JSONResponse(getBalance), nil
}
