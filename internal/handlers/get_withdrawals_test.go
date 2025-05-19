package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	"github.com/PiskarevSA/minimarket-points/internal/domain/entities"
	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
	"github.com/PiskarevSA/minimarket-points/internal/usecases"
)

type getWithdrawalsSuite struct {
	suite.Suite
	handlers    *Handlers
	mockUsecase *MockGetWithdrawalsUsecase
}

func TestGetWithdrawals(t *testing.T) {
	log.Logger = log.Logger.Output(io.Discard)

	suite.Run(t, new(getWithdrawalsSuite))
}

func (s *getWithdrawalsSuite) SetupTest() {
	s.mockUsecase = &MockGetWithdrawalsUsecase{}
	s.handlers = &Handlers{getWithdrawalsUsecase: s.mockUsecase}

	getJwtFromContext = func(
		_ context.Context,
		_ string,
	) (*jwt.Token, bool) {
		return &jwt.Token{}, true
	}

	getUserIdFromJwt = func(
		_ *jwt.Token,
		_ string,
	) (uuid.UUID, error) {
		return uuid.New(), nil
	}
}

func (s *getWithdrawalsSuite) TestGetWithdrawals_Success() {
	s.mockUsecase.DoFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		offset,
		limit *int32,
	) ([]entities.Transaction, error) {
		return nil, nil
	}

	response, err := s.handlers.GetWithdrawals(
		context.Background(),
		oapi.GetWithdrawalsRequestObject{},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitGetWithdrawalsResponse(rr)

	s.Require().Equal(http.StatusOK, rr.Result().StatusCode)
}

func (s *getWithdrawalsSuite) TestGetWithdrawals_NoJwt() {
	getJwtFromContext = func(
		_ context.Context,
		_ string,
	) (*jwt.Token, bool) {
		return nil, false
	}

	response, err := s.handlers.GetWithdrawals(
		context.Background(),
		oapi.GetWithdrawalsRequestObject{},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitGetWithdrawalsResponse(rr)

	s.Require().Equal(http.StatusInternalServerError, rr.Result().StatusCode)
}

func (s *getWithdrawalsSuite) TestGetWithdrawals_NoUserIdInClaims() {
	getUserIdFromJwt = func(
		_ *jwt.Token,
		_ string,
	) (uuid.UUID, error) {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	response, err := s.handlers.GetWithdrawals(
		context.Background(),
		oapi.GetWithdrawalsRequestObject{},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitGetWithdrawalsResponse(rr)

	s.Require().Equal(http.StatusInternalServerError, rr.Result().StatusCode)
}

func (s *getWithdrawalsSuite) TestGetWithdrawals_ValidationError() {
	s.mockUsecase.DoFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		offset,
		limit *int32,
	) ([]entities.Transaction, error) {
		return nil, &usecases.ValidationError{}
	}

	response, err := s.handlers.GetWithdrawals(
		context.Background(),
		oapi.GetWithdrawalsRequestObject{},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitGetWithdrawalsResponse(rr)

	s.Require().Equal(http.StatusBadRequest, rr.Result().StatusCode)
}
