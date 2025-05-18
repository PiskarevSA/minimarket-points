package handlers

import (
	"context"
	"errors"

	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
	"github.com/PiskarevSA/minimarket-points/internal/usecases"
)

var deposit500JSONResponse = oapi.Deposit500JSONResponse{
	Code:    oapi.S1394,
	Message: oapi.InternalError,
}

func (h *Handlers) Deposit(
	ctx context.Context,
	request oapi.DepositRequestObject,
) (oapi.DepositResponseObject, error) {
	const op = "deposit"

	token, ok := getJwtFromContext(ctx, op)
	if !ok {
		return deposit500JSONResponse, nil
	}

	userId, err := getUserIdFromJwt(token, op)
	if err != nil {
		return deposit500JSONResponse, nil
	}

	proccessedAt, err := h.depositUsecase.Do(
		ctx,
		userId,
		request.Body.OrderNumber,
		request.Body.Amount,
	)
	if err != nil {
		var validationErr *usecases.ValidationError
		if errors.As(err, &validationErr) {
			return oapi.Deposit400JSONResponse{
				Code:    oapi.ValidationErrorCode(validationErr.Code),
				Field:   validationErr.Field,
				Message: validationErr.Message,
			}, nil
		}

		var businessErr *usecases.BusinessError
		if errors.As(err, &businessErr) {
			return oapi.Deposit422JSONResponse{
				Code:    oapi.DomainErrorCode(businessErr.Code),
				Message: businessErr.Message,
			}, nil
		}

		return deposit500JSONResponse, nil
	}

	return oapi.Deposit200JSONResponse{
		ProccessedAt: proccessedAt,
	}, nil
}
