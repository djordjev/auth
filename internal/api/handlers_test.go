package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/domain/mocks"
	"github.com/djordjev/auth/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var signUpRequest = `
	{
		"username": "djvukovic",
		"password": "testee",
		"email": "djvukovic@gmail.com",
		"role": "admin"
	}
`

type SignUpTestSuite struct {
	suite.Suite
	request    *http.Request
	user       domain.User
	setup      domain.Setup
	logger     *slog.Logger
	mockDomain *mocks.Domain
}

func (suite *SignUpTestSuite) SetupTest() {
	requestPayload := strings.NewReader(signUpRequest)

	suite.user = domain.User{
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
	}

	suite.logger = utils.NewSilentLogger()
	suite.setup = domain.NewSetup(context.TODO(), suite.logger)

	req, err := http.NewRequestWithContext(context.TODO(), "POST", "/signup", requestPayload)
	if err != nil {
		suite.FailNow("failed to create test request")
		return
	}

	suite.request = utils.InjectLoggerIntoContext(req, suite.logger)
	suite.mockDomain = mocks.NewDomain(suite.T())
}

func (suite *SignUpTestSuite) TestSuccess() {
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()

	api := NewApi(utils.Config{}, mux, suite.mockDomain, suite.logger)

	newUser := domain.User{
		ID:       31,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
	}

	suite.mockDomain.EXPECT().SignUp(mock.Anything, suite.user).Return(newUser, nil)

	api.postSignup(rr, suite.request)

	suite.Require().Equal(rr.Code, http.StatusOK)

	var actualResponse SignUpResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResponse); err != nil {
		suite.T().FailNow()
	}

	suite.Require().Equal(actualResponse, SignUpResponse{
		ID:       newUser.ID,
		Username: newUser.Username,
		Password: newUser.Password,
		Role:     newUser.Role,
	})
}

func (suite *SignUpTestSuite) TestErrorUserAlreadyExists() {
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()

	api := NewApi(utils.Config{}, mux, suite.mockDomain, suite.logger)
	suite.mockDomain.EXPECT().SignUp(mock.Anything, suite.user).Return(domain.User{}, domain.ErrUserAlreadyExists)

	api.postSignup(rr, suite.request)

	suite.Require().Equal(rr.Code, http.StatusBadRequest)
}

func (suite *SignUpTestSuite) TestErrorInternal() {
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()

	api := NewApi(utils.Config{}, mux, suite.mockDomain, suite.logger)
	suite.mockDomain.EXPECT().SignUp(mock.Anything, suite.user).Return(domain.User{}, errors.New("err"))

	api.postSignup(rr, suite.request)

	suite.Require().Equal(rr.Code, http.StatusInternalServerError)
}

func TestSignUpTestSuite(t *testing.T) {
	suite.Run(t, new(SignUpTestSuite))
}
