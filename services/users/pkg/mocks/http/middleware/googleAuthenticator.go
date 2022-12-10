// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

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

// AddUserToContext provides a mock function with given fields: ctx, idToken
func (_m *GoogleAuthenticator) AddUserToContext(ctx context.Context, idToken string) (context.Context, error) {
	ret := _m.Called(ctx, idToken)

	var r0 context.Context
	if rf, ok := ret.Get(0).(func(context.Context, string) context.Context); ok {
		r0 = rf(ctx, idToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, idToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GoogleAuthenticator_AddUserToContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUserToContext'
type GoogleAuthenticator_AddUserToContext_Call struct {
	*mock.Call
}

// AddUserToContext is a helper method to define mock.On call
//   - ctx context.Context
//   - idToken string
func (_e *GoogleAuthenticator_Expecter) AddUserToContext(ctx interface{}, idToken interface{}) *GoogleAuthenticator_AddUserToContext_Call {
	return &GoogleAuthenticator_AddUserToContext_Call{Call: _e.mock.On("AddUserToContext", ctx, idToken)}
}

func (_c *GoogleAuthenticator_AddUserToContext_Call) Run(run func(ctx context.Context, idToken string)) *GoogleAuthenticator_AddUserToContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *GoogleAuthenticator_AddUserToContext_Call) Return(_a0 context.Context, _a1 error) *GoogleAuthenticator_AddUserToContext_Call {
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
