package handlers

import (
	"context"
	"errors"

	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
	"github.com/PiskarevSA/minimarket-points/internal/usecases"
)

var withdrawals500JSONResponse = oapi.GetWithdrawals500JSONResponse{
	Code:    oapi.S1394,
	Message: oapi.InternalError,
}

func (h *Handlers) GetWithdrawals(
	ctx context.Context,
	request oapi.GetWithdrawalsRequestObject,
) (oapi.GetWithdrawalsResponseObject, error) {
	const op = "deposit"

	token, ok := getJwtFromContext(ctx, op)
	if !ok {
		return withdrawals500JSONResponse, nil
	}

	userId, err := getUserIdFromJwt(token, op)
	if err != nil {
		return withdrawals500JSONResponse, nil
	}

	_, err = h.getWithdrawalsUsecase.Do(
		ctx,
		userId,
		request.Params.Limit,
		request.Params.Offset,
	)
	if err != nil {
		var validationErr *usecases.ValidationError
		if errors.As(err, &validationErr) {
			return oapi.GetWithdrawals400JSONResponse{
				Code:    oapi.ValidationErrorCode(validationErr.Code),
				Field:   validationErr.Field,
				Message: validationErr.Message,
			}, nil
		}

		return withdrawals500JSONResponse, nil
	}

	return nil, nil
}
