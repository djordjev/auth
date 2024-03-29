// Code generated by mockery v2.34.2. DO NOT EDIT.

package domain

import mock "github.com/stretchr/testify/mock"

// MockDomain is an autogenerated mock type for the Domain type
type MockDomain struct {
	mock.Mock
}

type MockDomain_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDomain) EXPECT() *MockDomain_Expecter {
	return &MockDomain_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: setup, user
func (_m *MockDomain) Delete(setup Setup, user User) (bool, error) {
	ret := _m.Called(setup, user)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(Setup, User) (bool, error)); ok {
		return rf(setup, user)
	}
	if rf, ok := ret.Get(0).(func(Setup, User) bool); ok {
		r0 = rf(setup, user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(Setup, User) error); ok {
		r1 = rf(setup, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDomain_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockDomain_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - setup Setup
//   - user User
func (_e *MockDomain_Expecter) Delete(setup interface{}, user interface{}) *MockDomain_Delete_Call {
	return &MockDomain_Delete_Call{Call: _e.mock.On("Delete", setup, user)}
}

func (_c *MockDomain_Delete_Call) Run(run func(setup Setup, user User)) *MockDomain_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Setup), args[1].(User))
	})
	return _c
}

func (_c *MockDomain_Delete_Call) Return(deleted bool, err error) *MockDomain_Delete_Call {
	_c.Call.Return(deleted, err)
	return _c
}

func (_c *MockDomain_Delete_Call) RunAndReturn(run func(Setup, User) (bool, error)) *MockDomain_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// LogIn provides a mock function with given fields: setup, user
func (_m *MockDomain) LogIn(setup Setup, user User) (User, string, error) {
	ret := _m.Called(setup, user)

	var r0 User
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(Setup, User) (User, string, error)); ok {
		return rf(setup, user)
	}
	if rf, ok := ret.Get(0).(func(Setup, User) User); ok {
		r0 = rf(setup, user)
	} else {
		r0 = ret.Get(0).(User)
	}

	if rf, ok := ret.Get(1).(func(Setup, User) string); ok {
		r1 = rf(setup, user)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(Setup, User) error); ok {
		r2 = rf(setup, user)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockDomain_LogIn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LogIn'
type MockDomain_LogIn_Call struct {
	*mock.Call
}

// LogIn is a helper method to define mock.On call
//   - setup Setup
//   - user User
func (_e *MockDomain_Expecter) LogIn(setup interface{}, user interface{}) *MockDomain_LogIn_Call {
	return &MockDomain_LogIn_Call{Call: _e.mock.On("LogIn", setup, user)}
}

func (_c *MockDomain_LogIn_Call) Run(run func(setup Setup, user User)) *MockDomain_LogIn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Setup), args[1].(User))
	})
	return _c
}

func (_c *MockDomain_LogIn_Call) Return(existing User, sessionKey string, err error) *MockDomain_LogIn_Call {
	_c.Call.Return(existing, sessionKey, err)
	return _c
}

func (_c *MockDomain_LogIn_Call) RunAndReturn(run func(Setup, User) (User, string, error)) *MockDomain_LogIn_Call {
	_c.Call.Return(run)
	return _c
}

// Logout provides a mock function with given fields: setup, token
func (_m *MockDomain) Logout(setup Setup, token string) error {
	ret := _m.Called(setup, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(Setup, string) error); ok {
		r0 = rf(setup, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDomain_Logout_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Logout'
type MockDomain_Logout_Call struct {
	*mock.Call
}

// Logout is a helper method to define mock.On call
//   - setup Setup
//   - token string
func (_e *MockDomain_Expecter) Logout(setup interface{}, token interface{}) *MockDomain_Logout_Call {
	return &MockDomain_Logout_Call{Call: _e.mock.On("Logout", setup, token)}
}

func (_c *MockDomain_Logout_Call) Run(run func(setup Setup, token string)) *MockDomain_Logout_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Setup), args[1].(string))
	})
	return _c
}

func (_c *MockDomain_Logout_Call) Return(err error) *MockDomain_Logout_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockDomain_Logout_Call) RunAndReturn(run func(Setup, string) error) *MockDomain_Logout_Call {
	_c.Call.Return(run)
	return _c
}

// ResetPasswordRequest provides a mock function with given fields: setup, user
func (_m *MockDomain) ResetPasswordRequest(setup Setup, user User) (User, error) {
	ret := _m.Called(setup, user)

	var r0 User
	var r1 error
	if rf, ok := ret.Get(0).(func(Setup, User) (User, error)); ok {
		return rf(setup, user)
	}
	if rf, ok := ret.Get(0).(func(Setup, User) User); ok {
		r0 = rf(setup, user)
	} else {
		r0 = ret.Get(0).(User)
	}

	if rf, ok := ret.Get(1).(func(Setup, User) error); ok {
		r1 = rf(setup, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDomain_ResetPasswordRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetPasswordRequest'
type MockDomain_ResetPasswordRequest_Call struct {
	*mock.Call
}

// ResetPasswordRequest is a helper method to define mock.On call
//   - setup Setup
//   - user User
func (_e *MockDomain_Expecter) ResetPasswordRequest(setup interface{}, user interface{}) *MockDomain_ResetPasswordRequest_Call {
	return &MockDomain_ResetPasswordRequest_Call{Call: _e.mock.On("ResetPasswordRequest", setup, user)}
}

func (_c *MockDomain_ResetPasswordRequest_Call) Run(run func(setup Setup, user User)) *MockDomain_ResetPasswordRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Setup), args[1].(User))
	})
	return _c
}

func (_c *MockDomain_ResetPasswordRequest_Call) Return(sentTo User, err error) *MockDomain_ResetPasswordRequest_Call {
	_c.Call.Return(sentTo, err)
	return _c
}

func (_c *MockDomain_ResetPasswordRequest_Call) RunAndReturn(run func(Setup, User) (User, error)) *MockDomain_ResetPasswordRequest_Call {
	_c.Call.Return(run)
	return _c
}

// Session provides a mock function with given fields: setup, token
func (_m *MockDomain) Session(setup Setup, token string) (User, error) {
	ret := _m.Called(setup, token)

	var r0 User
	var r1 error
	if rf, ok := ret.Get(0).(func(Setup, string) (User, error)); ok {
		return rf(setup, token)
	}
	if rf, ok := ret.Get(0).(func(Setup, string) User); ok {
		r0 = rf(setup, token)
	} else {
		r0 = ret.Get(0).(User)
	}

	if rf, ok := ret.Get(1).(func(Setup, string) error); ok {
		r1 = rf(setup, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDomain_Session_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Session'
type MockDomain_Session_Call struct {
	*mock.Call
}

// Session is a helper method to define mock.On call
//   - setup Setup
//   - token string
func (_e *MockDomain_Expecter) Session(setup interface{}, token interface{}) *MockDomain_Session_Call {
	return &MockDomain_Session_Call{Call: _e.mock.On("Session", setup, token)}
}

func (_c *MockDomain_Session_Call) Run(run func(setup Setup, token string)) *MockDomain_Session_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Setup), args[1].(string))
	})
	return _c
}

func (_c *MockDomain_Session_Call) Return(user User, err error) *MockDomain_Session_Call {
	_c.Call.Return(user, err)
	return _c
}

func (_c *MockDomain_Session_Call) RunAndReturn(run func(Setup, string) (User, error)) *MockDomain_Session_Call {
	_c.Call.Return(run)
	return _c
}

// SignUp provides a mock function with given fields: setup, user
func (_m *MockDomain) SignUp(setup Setup, user User) (User, error) {
	ret := _m.Called(setup, user)

	var r0 User
	var r1 error
	if rf, ok := ret.Get(0).(func(Setup, User) (User, error)); ok {
		return rf(setup, user)
	}
	if rf, ok := ret.Get(0).(func(Setup, User) User); ok {
		r0 = rf(setup, user)
	} else {
		r0 = ret.Get(0).(User)
	}

	if rf, ok := ret.Get(1).(func(Setup, User) error); ok {
		r1 = rf(setup, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDomain_SignUp_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SignUp'
type MockDomain_SignUp_Call struct {
	*mock.Call
}

// SignUp is a helper method to define mock.On call
//   - setup Setup
//   - user User
func (_e *MockDomain_Expecter) SignUp(setup interface{}, user interface{}) *MockDomain_SignUp_Call {
	return &MockDomain_SignUp_Call{Call: _e.mock.On("SignUp", setup, user)}
}

func (_c *MockDomain_SignUp_Call) Run(run func(setup Setup, user User)) *MockDomain_SignUp_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Setup), args[1].(User))
	})
	return _c
}

func (_c *MockDomain_SignUp_Call) Return(newUser User, err error) *MockDomain_SignUp_Call {
	_c.Call.Return(newUser, err)
	return _c
}

func (_c *MockDomain_SignUp_Call) RunAndReturn(run func(Setup, User) (User, error)) *MockDomain_SignUp_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyAccount provides a mock function with given fields: setup, token
func (_m *MockDomain) VerifyAccount(setup Setup, token string) (bool, error) {
	ret := _m.Called(setup, token)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(Setup, string) (bool, error)); ok {
		return rf(setup, token)
	}
	if rf, ok := ret.Get(0).(func(Setup, string) bool); ok {
		r0 = rf(setup, token)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(Setup, string) error); ok {
		r1 = rf(setup, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDomain_VerifyAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyAccount'
type MockDomain_VerifyAccount_Call struct {
	*mock.Call
}

// VerifyAccount is a helper method to define mock.On call
//   - setup Setup
//   - token string
func (_e *MockDomain_Expecter) VerifyAccount(setup interface{}, token interface{}) *MockDomain_VerifyAccount_Call {
	return &MockDomain_VerifyAccount_Call{Call: _e.mock.On("VerifyAccount", setup, token)}
}

func (_c *MockDomain_VerifyAccount_Call) Run(run func(setup Setup, token string)) *MockDomain_VerifyAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Setup), args[1].(string))
	})
	return _c
}

func (_c *MockDomain_VerifyAccount_Call) Return(verified bool, err error) *MockDomain_VerifyAccount_Call {
	_c.Call.Return(verified, err)
	return _c
}

func (_c *MockDomain_VerifyAccount_Call) RunAndReturn(run func(Setup, string) (bool, error)) *MockDomain_VerifyAccount_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyPasswordReset provides a mock function with given fields: setup, token, password
func (_m *MockDomain) VerifyPasswordReset(setup Setup, token string, password string) (User, error) {
	ret := _m.Called(setup, token, password)

	var r0 User
	var r1 error
	if rf, ok := ret.Get(0).(func(Setup, string, string) (User, error)); ok {
		return rf(setup, token, password)
	}
	if rf, ok := ret.Get(0).(func(Setup, string, string) User); ok {
		r0 = rf(setup, token, password)
	} else {
		r0 = ret.Get(0).(User)
	}

	if rf, ok := ret.Get(1).(func(Setup, string, string) error); ok {
		r1 = rf(setup, token, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDomain_VerifyPasswordReset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyPasswordReset'
type MockDomain_VerifyPasswordReset_Call struct {
	*mock.Call
}

// VerifyPasswordReset is a helper method to define mock.On call
//   - setup Setup
//   - token string
//   - password string
func (_e *MockDomain_Expecter) VerifyPasswordReset(setup interface{}, token interface{}, password interface{}) *MockDomain_VerifyPasswordReset_Call {
	return &MockDomain_VerifyPasswordReset_Call{Call: _e.mock.On("VerifyPasswordReset", setup, token, password)}
}

func (_c *MockDomain_VerifyPasswordReset_Call) Run(run func(setup Setup, token string, password string)) *MockDomain_VerifyPasswordReset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Setup), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDomain_VerifyPasswordReset_Call) Return(updated User, err error) *MockDomain_VerifyPasswordReset_Call {
	_c.Call.Return(updated, err)
	return _c
}

func (_c *MockDomain_VerifyPasswordReset_Call) RunAndReturn(run func(Setup, string, string) (User, error)) *MockDomain_VerifyPasswordReset_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDomain creates a new instance of MockDomain. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDomain(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDomain {
	mock := &MockDomain{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
