package usecases

import (
	"context"
	"time"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/google/uuid"
)

type AdjustBalanceRepo interface {
	AdjustBalance(
		ctx context.Context,
		userId uuid.UUID,
		orderNumber objects.OrderNumber,
		operation objects.Operation,
		amount objects.Amount,
		updatedAt time.Time,
	)
}
