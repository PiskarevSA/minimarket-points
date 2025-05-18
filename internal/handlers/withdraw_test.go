package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PiskarevSA/minimarket-points/internal/gen/oapi"
	"github.com/PiskarevSA/minimarket-points/internal/usecases"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type withdrawSuite struct {
	suite.Suite
	handlers    *Handlers
	mockUsecase *MockWithdrawUsecase
}

func TestWithdraw(t *testing.T) {
	log.Logger = log.Logger.Output(io.Discard)
	suite.Run(t, new(withdrawSuite))
}

func (s *withdrawSuite) SetupTest() {
	s.mockUsecase = &MockWithdrawUsecase{}
	s.handlers = &Handlers{withdrawUsecase: s.mockUsecase}

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

func (s *withdrawSuite) TestWithdraw_Success() {
	s.mockUsecase.DoFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		orderNumber,
		amount string,
	) (time.Time, error) {
		return time.Now(), nil
	}

	response, err := s.handlers.Withdraw(
		context.Background(),
		oapi.WithdrawRequestObject{
			Body: &oapi.WithdrawRequest{
				Amount:      "",
				OrderNumber: "",
			}},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitWithdrawResponse(rr)

	s.Require().Equal(rr.Result().StatusCode, http.StatusOK)
}

func (s *withdrawSuite) TestWithdraw_NoJwt() {
	getJwtFromContext = func(
		_ context.Context,
		_ string,
	) (*jwt.Token, bool) {
		return nil, false
	}

	response, err := s.handlers.Withdraw(
		context.Background(),
		oapi.WithdrawRequestObject{
			Body: &oapi.WithdrawRequest{
				Amount:      "",
				OrderNumber: "",
			}},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitWithdrawResponse(rr)

	s.Require().Equal(rr.Result().StatusCode, http.StatusInternalServerError)
}

func (s *withdrawSuite) TestWithdraw_NoUserIdInClaims() {
	getUserIdFromJwt = func(
		_ *jwt.Token,
		_ string,
	) (uuid.UUID, error) {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	response, err := s.handlers.Withdraw(
		context.Background(),
		oapi.WithdrawRequestObject{
			Body: &oapi.WithdrawRequest{
				Amount:      "",
				OrderNumber: "",
			}},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitWithdrawResponse(rr)

	s.Require().Equal(rr.Result().StatusCode, http.StatusInternalServerError)
}

func (s *withdrawSuite) TestWithdraw_ValidationError() {
	s.mockUsecase.DoFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		orderNumber,
		amount string,
	) (time.Time, error) {
		return time.Now(), &usecases.ValidationError{}
	}

	response, err := s.handlers.Withdraw(
		context.Background(),
		oapi.WithdrawRequestObject{
			Body: &oapi.WithdrawRequest{
				Amount:      "",
				OrderNumber: "",
			}},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitWithdrawResponse(rr)

	s.Require().Equal(rr.Result().StatusCode, http.StatusBadRequest)
}

func (s *withdrawSuite) TestUploadOrder_BusinessError() {
	s.mockUsecase.DoFunc = func(
		ctx context.Context,
		userId uuid.UUID,
		orderNumber,
		amount string,
	) (time.Time, error) {
		return time.Now(), &usecases.BusinessError{}
	}

	response, err := s.handlers.Withdraw(
		context.Background(),
		oapi.WithdrawRequestObject{
			Body: &oapi.WithdrawRequest{
				Amount:      "",
				OrderNumber: "",
			}},
	)

	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	response.VisitWithdrawResponse(rr)

	s.Require().Equal(rr.Result().StatusCode, http.StatusUnprocessableEntity)
}
