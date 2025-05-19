package dto

import (
	"github.com/PiskarevSA/minimarket-points/internal/domain/entities"
	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
	"github.com/PiskarevSA/minimarket-points/internal/gen/sqlc/postgresql"
	"github.com/PiskarevSA/minimarket-points/pkg/pgx/convtype"
)

func TransactionsToGetWithdrawalsResponse(
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

func GetTransactionRowsTransactions(rows []postgresql.Transaction) []entities.Transaction {
	result := make([]entities.Transaction, len(rows))

	for i, row := range rows {
		result[i].SetId(row.Id)
		result[i].SetUserId(row.UserId)

		orderNumber := objects.OrderNumber(row.OrderNumber)
		result[i].SetOrderNumber(orderNumber)

		operation := objects.Operation(row.Operation)
		result[i].SetOperation(operation)

		decimal, _ := convtype.NumericToDecimal(row.Amount)
		amount := objects.Amount(decimal)

		result[i].SetAmount(amount)
		result[i].SetProcessedAt(row.Proccessedat)
	}

	return result
}
