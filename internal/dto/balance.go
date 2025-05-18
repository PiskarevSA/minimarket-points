package dto

import (
	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
)

func BalanceToGetBalance(balance objects.Balance) oapi.GetBalanceResponse {
	return oapi.GetBalanceResponse{
		Available: balance.Available().String(),
		Withdrawn: balance.Withdrawn().String(),
		UpdatedAt: balance.UpdatedAt(),
	}
}
