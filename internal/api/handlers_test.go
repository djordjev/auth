package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djordjev/auth/internal/domain"
	"github.com/djordjev/auth/internal/domain/mocks"
	"github.com/djordjev/auth/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slog"
)

var signUpRequest = `
	{
		"username": "djvukovic",
		"password": "testee",
		"email": "djvukovic@gmail.com",
		"role": "admin"
	}
`

var logInRequest = `
	{
		"username": "djvukovic",
		"password": "testee",
		"email": "djvukovic@gmail.com"
	}
`

var deleteRequest = `
	{
		"email": "djvukovic@gmail.com",
		"password": "testeee"
	}
`

type SignUpTestSuite struct {
	suite.Suite
	rr             *httptest.ResponseRecorder
	api            *jsonApi
	logger         *slog.Logger
	domainExpector *mocks.Domain_Expecter
	requestBuilder func(string) *http.Request
}

func (suite *SignUpTestSuite) SetupTest() {
	suite.rr = httptest.NewRecorder()
	suite.logger = utils.NewSilentLogger()

	baseMock := mocks.NewDomain(suite.T())
	suite.domainExpector = baseMock.EXPECT()

	suite.api = NewApi(utils.Config{}, http.NewServeMux(), baseMock, suite.logger)
	suite.requestBuilder = utils.RequestBuilder("POST", "/signup")
}

func (suite *SignUpTestSuite) TestSignUpSuccess() {
	request := suite.requestBuilder(signUpRequest)

	newUser := domain.User{
		ID:       31,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
	}

	userMatcher := mock.MatchedBy(func(user domain.User) bool {
		return user.Username == "djvukovic" &&
			user.Password == "testee" &&
			user.Email == "djvukovic@gmail.com" &&
			user.Role == "admin"
	})

	suite.domainExpector.SignUp(mock.Anything, userMatcher).Return(newUser, nil)

	suite.api.postSignup(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusOK)

	var response SignUpResponse
	json.Unmarshal(suite.rr.Body.Bytes(), &response)

	suite.Require().Equal(response.ID, newUser.ID)
	suite.Require().Equal(response.Username, newUser.Username)
	suite.Require().Equal(response.Email, newUser.Email)
	suite.Require().Equal(response.Role, newUser.Role)

}

func (suite *SignUpTestSuite) TestSignUpErrorUserAlreadyExists() {
	request := suite.requestBuilder(signUpRequest)

	userMatcher := mock.MatchedBy(func(user domain.User) bool {
		return user.Username == "djvukovic" &&
			user.Password == "testee" &&
			user.Email == "djvukovic@gmail.com" &&
			user.Role == "admin"
	})

	suite.domainExpector.SignUp(mock.Anything, userMatcher).Return(domain.User{}, domain.ErrUserAlreadyExists)

	suite.api.postSignup(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(
		suite.rr.Body.String(),
		utils.ErrorJSON("user with email djvukovic@gmail.com already exists"),
	)
}

func (suite *SignUpTestSuite) TestSignUpErrorInternal() {
	request := suite.requestBuilder(signUpRequest)

	userMatcher := mock.MatchedBy(func(user domain.User) bool {
		return user.Username == "djvukovic" &&
			user.Password == "testee" &&
			user.Email == "djvukovic@gmail.com" &&
			user.Role == "admin"
	})

	suite.domainExpector.SignUp(mock.Anything, userMatcher).Return(domain.User{}, errors.New("err"))

	suite.api.postSignup(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusInternalServerError)
	suite.Require().JSONEq(suite.rr.Body.String(), utils.ErrorJSON("internal server error"))
}

func (suite *SignUpTestSuite) TestSignUpValidationFail() {
	request := suite.requestBuilder(`
		{
			"username": "djvukovic",
			"password": "re",
			"email": "djvukovic@gmail.com",
			"role": "admin"
		}
	`)

	suite.api.postSignup(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(
		suite.rr.Body.String(),
		utils.ErrorJSON("password must have at least 5 characters"),
	)
}

type LogInTestSuite struct {
	suite.Suite
	rr             *httptest.ResponseRecorder
	api            *jsonApi
	logger         *slog.Logger
	domainExpector *mocks.Domain_Expecter
	requestBuilder func(string) *http.Request
}

func (suite *LogInTestSuite) SetupTest() {
	suite.rr = httptest.NewRecorder()
	suite.logger = utils.NewSilentLogger()

	baseMock := mocks.NewDomain(suite.T())
	suite.domainExpector = baseMock.EXPECT()

	suite.api = NewApi(utils.Config{}, http.NewServeMux(), baseMock, suite.logger)
	suite.requestBuilder = utils.RequestBuilder("POST", "/login")
}

func (suite *LogInTestSuite) TestLoginValidationFail() {
	invalidReq := `
		{
			"username": "djvukovic",
			"password": "tst",
			"email": "djvukovic@gmail.com"
		}
	`
	request := suite.requestBuilder(invalidReq)

	suite.api.postLogin(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(suite.rr.Body.String(), `{"error": "incorrect password" }`)
}

func (suite *LogInTestSuite) TestLoginBadRequest() {
	request := suite.requestBuilder(`{"something"}`)
	suite.api.postLogin(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(suite.rr.Body.String(), `{"error": "bad request"}`)
}

func (suite *LogInTestSuite) TestLoginSuccess() {
	user := domain.User{
		ID:       23121,
		Email:    "djvukovic@gmail.com",
		Username: "djvukovic",
		Password: "testee",
		Role:     "admin",
		Verified: true,
	}

	suite.domainExpector.LogIn(mock.Anything, mock.Anything).Return(user, nil)
	request := suite.requestBuilder(logInRequest)

	suite.api.postLogin(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusOK)

	var response LogInResponse
	json.Unmarshal(suite.rr.Body.Bytes(), &response)

	suite.Require().Equal(response.Email, user.Email)
	suite.Require().Equal(response.ID, user.ID)
	suite.Require().Equal(response.Role, user.Role)
	suite.Require().Equal(response.Verified, user.Verified)
	suite.Require().Equal(response.Username, user.Username)
}

func (suite *LogInTestSuite) TestLoginInvalidCredentials() {
	suite.domainExpector.LogIn(mock.Anything, mock.Anything).Return(domain.User{}, domain.ErrInvalidCredentials)
	request := suite.requestBuilder(logInRequest)

	suite.api.postLogin(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(suite.rr.Body.String(), utils.ErrorJSON("invalid credentials"))
}

func (suite *LogInTestSuite) TestLoginOtherError() {
	suite.domainExpector.LogIn(mock.Anything, mock.Anything).Return(domain.User{}, errors.New("random"))
	request := suite.requestBuilder(logInRequest)

	suite.api.postLogin(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(suite.rr.Body.String(), utils.ErrorJSON("failed login attempt"))
}

type DeleteTestSuite struct {
	suite.Suite
	rr             *httptest.ResponseRecorder
	api            *jsonApi
	logger         *slog.Logger
	domainExpector *mocks.Domain_Expecter
	requestBuilder func(string) *http.Request
}

func (suite *DeleteTestSuite) SetupTest() {
	suite.rr = httptest.NewRecorder()
	suite.logger = utils.NewSilentLogger()

	baseMock := mocks.NewDomain(suite.T())
	suite.domainExpector = baseMock.EXPECT()

	suite.api = NewApi(utils.Config{}, http.NewServeMux(), baseMock, suite.logger)
	suite.requestBuilder = utils.RequestBuilder("DELETE", "/account")
}

func (suite *DeleteTestSuite) TestDeleteSuccess() {
	request := suite.requestBuilder(deleteRequest)
	suite.domainExpector.Delete(mock.Anything, mock.Anything).Return(true, nil)

	suite.api.deleteAccount(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusOK)
	var response DeleteAccountResponse
	json.Unmarshal(suite.rr.Body.Bytes(), &response)
	suite.Require().Equal(response.Success, true)

}

func (suite *DeleteTestSuite) TestDeleteValidationFailed() {
	request := suite.requestBuilder(`{"email": "something"}`)

	suite.api.deleteAccount(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(suite.rr.Body.String(), utils.ErrorJSON("missing password"))
}

func (suite *DeleteTestSuite) TestDeleteUserNotExist() {
	request := suite.requestBuilder(deleteRequest)

	suite.domainExpector.Delete(mock.Anything, mock.Anything).Return(false, domain.ErrUserNotExist)

	suite.api.deleteAccount(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(suite.rr.Body.String(), utils.ErrorJSON("user does not exists"))
}

func (suite *DeleteTestSuite) TestDeleteAuthenticationFailed() {
	request := suite.requestBuilder(deleteRequest)

	suite.domainExpector.Delete(mock.Anything, mock.Anything).Return(false, domain.ErrInvalidCredentials)

	suite.api.deleteAccount(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusBadRequest)
	suite.Require().JSONEq(suite.rr.Body.String(), utils.ErrorJSON("authentication failed"))
}

func (suite *DeleteTestSuite) TestDeleteError() {
	request := suite.requestBuilder(deleteRequest)

	suite.domainExpector.Delete(mock.Anything, mock.Anything).Return(false, errors.New("something"))

	suite.api.deleteAccount(suite.rr, request)

	suite.Require().Equal(suite.rr.Code, http.StatusInternalServerError)
	suite.Require().JSONEq(suite.rr.Body.String(), utils.ErrorJSON("internal server error"))
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(SignUpTestSuite))
	suite.Run(t, new(LogInTestSuite))
	suite.Run(t, new(DeleteTestSuite))
}
