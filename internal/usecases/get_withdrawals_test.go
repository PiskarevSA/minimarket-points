package usecases

import (
	"context"
	"errors"
	"io"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/PiskarevSA/minimarket-points/internal/domain/entities"
	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
	"github.com/PiskarevSA/minimarket-points/pkg/strgen"
)

type getWithdrawalsSuite struct {
	suite.Suite
	usecase     *GetWithdrawals
	mockStorage *MockGetWithdrawalsRepo

	rawOffset int32
	rawLimit  int32
}

func TestGetWithdrawals(t *testing.T) {
	log.Logger = log.Logger.Output(io.Discard)

	suite.Run(t, new(getWithdrawalsSuite))
}

func (s *getWithdrawalsSuite) SetupTest() {
	s.mockStorage = &MockGetWithdrawalsRepo{}
	limitMax := int32(35)
	s.usecase = NewGetWithdrawals(limitMax, s.mockStorage)

	s.rawOffset = 0
	s.rawLimit = 15
}

func (s *getWithdrawalsSuite) generateTxs(n int) []entities.Transaction {
	txs := make([]entities.Transaction, n)

	operations := []objects.Operation{
		objects.OperationDeposit,
		objects.OperationWithdraw,
	}
	numOperations := len(operations)

	onlyNumbers := strgen.WithAlphabet(strgen.AlphabetOnlyNumbers)
	strgen := strgen.New(onlyNumbers)

	for i := range n {
		txs[i].SetId(rand.Int32())
		txs[i].SetUserId(uuid.New())

		orderNumber := objects.OrderNumber(strgen.Generate(10))
		txs[i].SetOrderNumber(orderNumber)

		operation := operations[rand.IntN(numOperations)]
		txs[i].SetOperation(operation)

		decimal := decimal.NewFromFloat(rand.Float64() * 100)
		txs[i].SetAmount(objects.Amount(decimal))

		proccessedAt := time.Now().
			Add(-time.Duration(rand.IntN(1000)) * time.Hour)
		txs[i].SetProcessedAt(proccessedAt)
	}

	return txs
}

func (s *getWithdrawalsSuite) TestGetWithdrawals_Success() {
	expectedTxs := s.generateTxs(rand.IntN(10))
	s.mockStorage.GetWithdrawalsFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		offset,
		limit int32,
	) ([]entities.Transaction, error) {
		return expectedTxs, nil
	}

	txs, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		&s.rawOffset,
		&s.rawLimit,
	)

	s.Require().NoError(err)
	s.Require().ElementsMatch(expectedTxs, txs)
}

func (s *getWithdrawalsSuite) TestGetWithdrawals_InvalidPagination() {
	rawOffset := int32(15)
	rawLimit := int32(50)

	txs, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		&rawOffset,
		&rawLimit,
	)

	var validationErr *ValidationError

	s.Require().ErrorAs(err, &validationErr)
	s.Require().Equal("V1019", validationErr.Code)

	s.Require().Empty(txs)
}

func (s *getWithdrawalsSuite) TestGetWithdrawals_NoTransactions() {
	s.mockStorage.GetWithdrawalsFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		offset,
		limit int32,
	) ([]entities.Transaction, error) {
		return nil, repo.ErrNoTransactionsFound
	}

	txs, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		&s.rawOffset,
		&s.rawLimit,
	)

	s.Require().NoError(err)
	s.Require().Empty(txs)
}

func (s *getWithdrawalsSuite) TestGetWithdrawals_UnexpectedError() {
	errUnexpectedError := errors.New("unexpected error")

	s.mockStorage.GetWithdrawalsFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		offset,
		limit int32,
	) ([]entities.Transaction, error) {
		return nil, errUnexpectedError
	}

	txs, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		&s.rawOffset,
		&s.rawLimit,
	)

	s.Require().ErrorIs(err, errUnexpectedError)
	s.Require().Empty(txs)
}
