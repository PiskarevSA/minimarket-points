package handlers

import (
	"context"
	"errors"

	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
	"github.com/PiskarevSA/minimarket-points/internal/usecases"
)

var withdraw500JSONResponse = oapi.Withdraw500JSONResponse{
	Code:    oapi.S1394,
	Message: oapi.InternalError,
}

func (h *Handlers) Withdraw(
	ctx context.Context,
	request oapi.WithdrawRequestObject,
) (oapi.WithdrawResponseObject, error) {
	const op = "deposit"

	token, ok := getJwtFromContext(ctx, op)
	if !ok {
		return withdraw500JSONResponse, nil
	}

	userId, err := getUserIdFromJwt(token, op)
	if err != nil {
		return withdraw500JSONResponse, nil
	}

	proccessedAt, err := h.withdrawUsecase.Do(
		ctx,
		userId,
		request.Body.OrderNumber,
		request.Body.Amount,
	)
	if err != nil {
		var validationErr *usecases.ValidationError
		if errors.As(err, &validationErr) {
			return oapi.Withdraw400JSONResponse{
				Code:    oapi.ValidationErrorCode(validationErr.Code),
				Field:   validationErr.Field,
				Message: validationErr.Message,
			}, nil
		}

		var businessErr *usecases.BusinessError
		if errors.As(err, &businessErr) {
			return oapi.Withdraw422JSONResponse{
				Code:    oapi.DomainErrorCode(businessErr.Code),
				Message: businessErr.Message,
			}, nil
		}

		return withdraw500JSONResponse, nil
	}

	return oapi.Withdraw200JSONResponse{
		ProccessedAt: proccessedAt,
	}, nil
}
