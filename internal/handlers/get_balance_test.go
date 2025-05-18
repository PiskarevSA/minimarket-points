package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PiskarevSA/minimarket-points/internal/domain/objects"
	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type getBalanceSuite struct {
	suite.Suite
	handlers    *Handlers
	mockUsecase *MockGetBalanceUsecase
}

func TestGetBalance(t *testing.T) {
	log.Logger = log.Logger.Output(io.Discard)
	suite.Run(t, new(getBalanceSuite))
}

func (s *getBalanceSuite) SetupTest() {
	s.mockUsecase = &MockGetBalanceUsecase{}
	s.handlers = &Handlers{getBalanceUsecase: s.mockUsecase}

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

func (s *getBalanceSuite) TestGetBalance_Success() {
	s.mockUsecase.DoFunc = func(
		ctx context.Context,
		userId uuid.UUID,
	) (objects.Balance, error) {
		return objects.NilBalance, nil
	}

	response, err := s.handlers.GetBalance(
		context.Background(),
		oapi.GetBalanceRequestObject{},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitGetBalanceResponse(rr)

	s.Require().Equal(rr.Result().StatusCode, http.StatusOK)
}

func (s *getBalanceSuite) TestGetBalance_NoJwt() {
	getJwtFromContext = func(
		_ context.Context,
		_ string,
	) (*jwt.Token, bool) {
		return nil, false
	}

	response, err := s.handlers.GetBalance(
		context.Background(),
		oapi.GetBalanceRequestObject{},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitGetBalanceResponse(rr)

	s.Require().Equal(rr.Result().StatusCode, http.StatusInternalServerError)
}

func (s *getBalanceSuite) TestGetBalance_NoUserIdInClaims() {
	getUserIdFromJwt = func(
		_ *jwt.Token,
		_ string,
	) (uuid.UUID, error) {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	response, err := s.handlers.GetBalance(
		context.Background(),
		oapi.GetBalanceRequestObject{},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitGetBalanceResponse(rr)

	s.Require().Equal(rr.Result().StatusCode, http.StatusInternalServerError)
}
