package dto

import (
	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
	"github.com/PiskarevSA/minimarket-points/internal/gen/sqlc/postgresql"
	"github.com/PiskarevSA/minimarket-points/pkg/pgx/convtype"
)

func BalanceToGetBalanceResponse(balance objects.Balance) oapi.GetBalanceResponse {
	return oapi.GetBalanceResponse{
		Available: balance.Available().String(),
		Withdrawn: balance.Withdrawn().String(),
		UpdatedAt: balance.UpdatedAt(),
	}
}

func GetBalanceRowToBalance(row postgresql.GetBalanceRow) objects.Balance {
	balance := objects.Balance{}

	available, _ := convtype.NumericToDecimal(row.Available)
	balance.SetAvailable(objects.Amount(available))

	withdrawn, _ := convtype.NumericToDecimal(row.Withdrawn)
	balance.SetWithdrawn(objects.Amount(withdrawn))

	balance.SetUpdatedAt(row.UpdatedAt)

	return balance
}
