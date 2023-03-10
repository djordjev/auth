// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/djordjev/auth/internal/domain/types"
)

// Domain is an autogenerated mock type for the Domain type
type Domain struct {
	mock.Mock
}

type Domain_Expecter struct {
	mock *mock.Mock
}

func (_m *Domain) EXPECT() *Domain_Expecter {
	return &Domain_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, user
func (_m *Domain) Delete(ctx context.Context, user types.User) (bool, error) {
	ret := _m.Called(ctx, user)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.User) (bool, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.User) bool); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Domain_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type Domain_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - user types.User
func (_e *Domain_Expecter) Delete(ctx interface{}, user interface{}) *Domain_Delete_Call {
	return &Domain_Delete_Call{Call: _e.mock.On("Delete", ctx, user)}
}

func (_c *Domain_Delete_Call) Run(run func(ctx context.Context, user types.User)) *Domain_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.User))
	})
	return _c
}

func (_c *Domain_Delete_Call) Return(deleted bool, err error) *Domain_Delete_Call {
	_c.Call.Return(deleted, err)
	return _c
}

func (_c *Domain_Delete_Call) RunAndReturn(run func(context.Context, types.User) (bool, error)) *Domain_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// LogIn provides a mock function with given fields: ctx, user
func (_m *Domain) LogIn(ctx context.Context, user types.User) (types.User, error) {
	ret := _m.Called(ctx, user)

	var r0 types.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.User) (types.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.User) types.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(types.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Domain_LogIn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LogIn'
type Domain_LogIn_Call struct {
	*mock.Call
}

// LogIn is a helper method to define mock.On call
//   - ctx context.Context
//   - user types.User
func (_e *Domain_Expecter) LogIn(ctx interface{}, user interface{}) *Domain_LogIn_Call {
	return &Domain_LogIn_Call{Call: _e.mock.On("LogIn", ctx, user)}
}

func (_c *Domain_LogIn_Call) Run(run func(ctx context.Context, user types.User)) *Domain_LogIn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.User))
	})
	return _c
}

func (_c *Domain_LogIn_Call) Return(exisingUser types.User, err error) *Domain_LogIn_Call {
	_c.Call.Return(exisingUser, err)
	return _c
}

func (_c *Domain_LogIn_Call) RunAndReturn(run func(context.Context, types.User) (types.User, error)) *Domain_LogIn_Call {
	_c.Call.Return(run)
	return _c
}

// LogOut provides a mock function with given fields: ctx, user
func (_m *Domain) LogOut(ctx context.Context, user types.User) (bool, error) {
	ret := _m.Called(ctx, user)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.User) (bool, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.User) bool); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Domain_LogOut_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LogOut'
type Domain_LogOut_Call struct {
	*mock.Call
}

// LogOut is a helper method to define mock.On call
//   - ctx context.Context
//   - user types.User
func (_e *Domain_Expecter) LogOut(ctx interface{}, user interface{}) *Domain_LogOut_Call {
	return &Domain_LogOut_Call{Call: _e.mock.On("LogOut", ctx, user)}
}

func (_c *Domain_LogOut_Call) Run(run func(ctx context.Context, user types.User)) *Domain_LogOut_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.User))
	})
	return _c
}

func (_c *Domain_LogOut_Call) Return(loggedOut bool, err error) *Domain_LogOut_Call {
	_c.Call.Return(loggedOut, err)
	return _c
}

func (_c *Domain_LogOut_Call) RunAndReturn(run func(context.Context, types.User) (bool, error)) *Domain_LogOut_Call {
	_c.Call.Return(run)
	return _c
}

// SignUp provides a mock function with given fields: ctx, user
func (_m *Domain) SignUp(ctx context.Context, user types.User) (types.User, error) {
	ret := _m.Called(ctx, user)

	var r0 types.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.User) (types.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.User) types.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(types.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Domain_SignUp_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SignUp'
type Domain_SignUp_Call struct {
	*mock.Call
}

// SignUp is a helper method to define mock.On call
//   - ctx context.Context
//   - user types.User
func (_e *Domain_Expecter) SignUp(ctx interface{}, user interface{}) *Domain_SignUp_Call {
	return &Domain_SignUp_Call{Call: _e.mock.On("SignUp", ctx, user)}
}

func (_c *Domain_SignUp_Call) Run(run func(ctx context.Context, user types.User)) *Domain_SignUp_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.User))
	})
	return _c
}

func (_c *Domain_SignUp_Call) Return(newUser types.User, err error) *Domain_SignUp_Call {
	_c.Call.Return(newUser, err)
	return _c
}

func (_c *Domain_SignUp_Call) RunAndReturn(run func(context.Context, types.User) (types.User, error)) *Domain_SignUp_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewDomain interface {
	mock.TestingT
	Cleanup(func())
}

// NewDomain creates a new instance of Domain. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDomain(t mockConstructorTestingTNewDomain) *Domain {
	mock := &Domain{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
