package dto

import (
	"github.com/PiskarevSA/minimarket-points/internal/domain/entities"
	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
)

func TransactionsToGetWithdrawals(
	txs []entities.Transaction,
) oapi.GetWithdrawalsResponse {
	result := make(oapi.GetWithdrawalsResponse, len(txs))

	for i, tx := range txs {
		result[i].Amount = tx.Amount().String()
		result[i].Id = tx.Id()
		result[i].OrderNumber = tx.OrderNumber().String()
		result[i].ProccessedAt = tx.ProccessedAt()
	}

	return result
}
