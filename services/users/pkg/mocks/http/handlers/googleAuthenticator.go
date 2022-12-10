// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	google "github.com/clubrizer/services/users/internal/authenticator/google"

	mock "github.com/stretchr/testify/mock"
)

// GoogleAuthenticator is an autogenerated mock type for the googleAuthenticator type
type GoogleAuthenticator struct {
	mock.Mock
}

type GoogleAuthenticator_Expecter struct {
	mock *mock.Mock
}

func (_m *GoogleAuthenticator) EXPECT() *GoogleAuthenticator_Expecter {
	return &GoogleAuthenticator_Expecter{mock: &_m.Mock}
}

// GetUserFromContext provides a mock function with given fields: ctx
func (_m *GoogleAuthenticator) GetUserFromContext(ctx context.Context) (*google.User, bool) {
	ret := _m.Called(ctx)

	var r0 *google.User
	if rf, ok := ret.Get(0).(func(context.Context) *google.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*google.User)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context) bool); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GoogleAuthenticator_GetUserFromContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserFromContext'
type GoogleAuthenticator_GetUserFromContext_Call struct {
	*mock.Call
}

// GetUserFromContext is a helper method to define mock.On call
//   - ctx context.Context
func (_e *GoogleAuthenticator_Expecter) GetUserFromContext(ctx interface{}) *GoogleAuthenticator_GetUserFromContext_Call {
	return &GoogleAuthenticator_GetUserFromContext_Call{Call: _e.mock.On("GetUserFromContext", ctx)}
}

func (_c *GoogleAuthenticator_GetUserFromContext_Call) Run(run func(ctx context.Context)) *GoogleAuthenticator_GetUserFromContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *GoogleAuthenticator_GetUserFromContext_Call) Return(_a0 *google.User, _a1 bool) *GoogleAuthenticator_GetUserFromContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewGoogleAuthenticator interface {
	mock.TestingT
	Cleanup(func())
}

// NewGoogleAuthenticator creates a new instance of GoogleAuthenticator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGoogleAuthenticator(t mockConstructorTestingTNewGoogleAuthenticator) *GoogleAuthenticator {
	mock := &GoogleAuthenticator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}