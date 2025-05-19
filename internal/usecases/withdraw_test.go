package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/repo"
	"github.com/PiskarevSA/minimarket-points/pkg/damm/damm"
	"github.com/PiskarevSA/minimarket-points/pkg/strgen"
)

type withdrawSuite struct {
	suite.Suite
	usecase     *Withdraw
	mockStorage *MockWithdrawRepo
	transactor  *MockTransactor
	strgen      *strgen.Generator

	errWhenCreatingTransaction error
	errWhenAdjustingBalance    error
}

func TestWithdraw(t *testing.T) {
	log.Logger = log.Logger.Output(io.Discard)

	suite.Run(t, new(withdrawSuite))
}

func (s *withdrawSuite) SetupTest() {
	s.mockStorage = &MockWithdrawRepo{}
	s.transactor = &MockTransactor{}
	s.usecase = NewWithdraw(s.mockStorage, s.transactor)

	onlyNumbers := strgen.WithAlphabet(strgen.AlphabetOnlyNumbers)
	s.strgen = strgen.New(onlyNumbers)

	s.errWhenCreatingTransaction = nil
	s.errWhenAdjustingBalance = nil

	s.transactor.TransactFunc = func(
		ctx context.Context,
		opts pgx.TxOptions,
		fn func(ctx context.Context) error,
	) error {
		return fn(ctx)
	}

	s.mockStorage.CreateTransactionInTxFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		orderNumber objects.OrderNumber,
		operation objects.Operation,
		amount objects.Amount,
		timestamp time.Time,
	) error {
		return s.errWhenCreatingTransaction
	}

	s.mockStorage.AdjustBalanceInTxFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		orderNumber objects.OrderNumber,
		operation objects.Operation,
		amount objects.Amount,
		updatedAt time.Time,
	) error {
		return s.errWhenAdjustingBalance
	}
}

func (s *withdrawSuite) TestWithdraw_Success() {
	now := time.Now()
	timeNow = func() time.Time { return now }

	orderNumber, err := damm.Append(s.strgen.Generate(9))
	s.Require().NoError(err)

	amount := fmt.Sprintf("%.2f", rand.Float64()*100)

	proccessedAt, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		orderNumber,
		amount,
	)

	s.Require().NoError(err)
	s.Require().Equal(now, proccessedAt)
}

func (s *withdrawSuite) TestWithdraw_InvalidOrderNumber() {
	invalidOrderNumber, err := damm.Append(s.strgen.Generate(12))
	s.Require().NoError(err)

	amount := fmt.Sprintf("%.2f", rand.Float64()*100)

	proccessedAt, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		invalidOrderNumber,
		amount,
	)

	var validationErr *ValidationError
	s.Require().ErrorAs(err, &validationErr)
	s.Require().Equal("V1107", validationErr.Code)

	s.Require().Equal(time.Time{}, proccessedAt)
}

func (s *withdrawSuite) TestWithdraw_InvalidAmount() {
	rawOrderNumber, err := damm.Append(s.strgen.Generate(9))
	s.Require().NoError(err)

	invalidAmount := "-1"

	proccessedAt, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		rawOrderNumber,
		invalidAmount,
	)

	var validationErr *ValidationError
	s.Require().ErrorAs(err, &validationErr)
	s.Require().Equal("V1142", validationErr.Code)

	s.Require().Equal(time.Time{}, proccessedAt)
}

func (s *withdrawSuite) TestWithdraw_InsufficientBalance() {
	s.errWhenAdjustingBalance = repo.ErrNotEnoughtBalance

	orderNumber, err := damm.Append(s.strgen.Generate(9))
	s.Require().NoError(err)

	amount := fmt.Sprintf("%.2f", rand.Float64()*100)

	proccessedAt, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		orderNumber,
		amount,
	)

	var businessErr *BusinessError

	s.Require().ErrorAs(err, &businessErr)
	s.Require().Equal("D1215", businessErr.Code)

	s.Require().Equal(time.Time{}, proccessedAt)
}

func (s *withdrawSuite) TestGetWithdrawals_UnexpectedError() {
	errUnexpectedError := errors.New("unexpected error")
	s.errWhenAdjustingBalance = errUnexpectedError

	orderNumber, err := damm.Append(s.strgen.Generate(9))
	s.Require().NoError(err)

	amount := fmt.Sprintf("%.2f", rand.Float64()*100)

	proccessedAt, err := s.usecase.Do(
		context.Background(),
		uuid.New(),
		orderNumber,
		amount,
	)

	s.Require().ErrorIs(err, errUnexpectedError)
	s.Require().Equal(time.Time{}, proccessedAt)
}
