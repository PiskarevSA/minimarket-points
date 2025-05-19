package usecases

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
)

type getBalanceSuite struct {
	suite.Suite
	usecase     *GetBalance
	mockStorage *MockGetBalanceRepo
}

func TestGetBalance(t *testing.T) {
	log.Logger = log.Logger.Output(io.Discard)

	suite.Run(t, new(getBalanceSuite))
}

func (s *getBalanceSuite) SetupTest() {
	s.mockStorage = &MockGetBalanceRepo{}
	s.usecase = NewGetBalance(s.mockStorage)
}

func (s *getBalanceSuite) TestGetBalance_Success() {
	dec := decimal.NewFromFloat(431.1)
	available := objects.Amount(dec)

	dec = decimal.NewFromFloat(212.1)
	withdrawn := objects.Amount(dec)

	updatedAt := time.Now()

	s.mockStorage.GetBalanceFunc = func(
		ctx context.Context,
		userId uuid.UUID,
	) (objects.Balance, error) {
		balance := objects.Balance{}

		balance.SetAvailable(available)
		balance.SetWithdrawn(withdrawn)
		balance.SetUpdatedAt(updatedAt)

		return balance, nil
	}

	balance, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
	)

	s.Require().NoError(err)
	s.Require().Equal(balance.Available(), available)
	s.Require().Equal(balance.Withdrawn(), withdrawn)
	s.Require().Equal(balance.UpdatedAt(), updatedAt)
}

func (s *getBalanceSuite) TestGetBalance_NoBalance() {
	now := time.Now()
	timeNow = func() time.Time { return now }

	s.mockStorage.GetBalanceFunc = func(
		ctx context.Context,
		userId uuid.UUID,
	) (objects.Balance, error) {
		return objects.NilBalance, repo.ErrNoBalanceFound
	}

	balance, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
	)

	s.Require().ErrorIs(err, repo.ErrNoBalanceFound)

	s.Require().Empty(balance.Available())
	s.Require().Empty(balance.Withdrawn())
	s.Require().Equal(now, balance.UpdatedAt())
}

func (s *getBalanceSuite) TestGetWithdrawals_UnexpectedError() {
	errUnexpectedError := errors.New("unexpected error")

	s.mockStorage.GetBalanceFunc = func(
		ctx context.Context,
		userId uuid.UUID,
	) (objects.Balance, error) {
		return objects.NilBalance, errUnexpectedError
	}

	balance, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
	)

	s.Require().ErrorIs(err, errUnexpectedError)
	s.Require().Empty(balance)
}
