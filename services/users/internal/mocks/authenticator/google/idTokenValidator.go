// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	idtoken "google.golang.org/api/idtoken"

	mock "github.com/stretchr/testify/mock"
)

// IdTokenValidator is an autogenerated mock type for the idTokenValidator type
type IdTokenValidator struct {
	mock.Mock
}

type IdTokenValidator_Expecter struct {
	mock *mock.Mock
}

func (_m *IdTokenValidator) EXPECT() *IdTokenValidator_Expecter {
	return &IdTokenValidator_Expecter{mock: &_m.Mock}
}

// Validate provides a mock function with given fields: ctx, idToken, audience
func (_m *IdTokenValidator) Validate(ctx context.Context, idToken string, audience string) (*idtoken.Payload, error) {
	ret := _m.Called(ctx, idToken, audience)

	var r0 *idtoken.Payload
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *idtoken.Payload); ok {
		r0 = rf(ctx, idToken, audience)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*idtoken.Payload)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, idToken, audience)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IdTokenValidator_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type IdTokenValidator_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - ctx context.Context
//   - idToken string
//   - audience string
func (_e *IdTokenValidator_Expecter) Validate(ctx interface{}, idToken interface{}, audience interface{}) *IdTokenValidator_Validate_Call {
	return &IdTokenValidator_Validate_Call{Call: _e.mock.On("Validate", ctx, idToken, audience)}
}

func (_c *IdTokenValidator_Validate_Call) Run(run func(ctx context.Context, idToken string, audience string)) *IdTokenValidator_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *IdTokenValidator_Validate_Call) Return(_a0 *idtoken.Payload, _a1 error) *IdTokenValidator_Validate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewIdTokenValidator interface {
	mock.TestingT
	Cleanup(func())
}

// NewIdTokenValidator creates a new instance of IdTokenValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIdTokenValidator(t mockConstructorTestingTNewIdTokenValidator) *IdTokenValidator {
	mock := &IdTokenValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}