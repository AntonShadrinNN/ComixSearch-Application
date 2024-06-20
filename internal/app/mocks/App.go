// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// App is an autogenerated mock type for the App type
type App struct {
	mock.Mock
}

// Search provides a mock function with given fields: ctx, keywords, limit
func (_m *App) Search(ctx context.Context, keywords []string, limit int) (map[string]string, error) {
	ret := _m.Called(ctx, keywords, limit)

	if len(ret) == 0 {
		panic("no return value specified for Search")
	}

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, int) (map[string]string, error)); ok {
		return rf(ctx, keywords, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string, int) map[string]string); ok {
		r0 = rf(ctx, keywords, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string, int) error); ok {
		r1 = rf(ctx, keywords, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewApp creates a new instance of App. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApp(t interface {
	mock.TestingT
	Cleanup(func())
}) *App {
	mock := &App{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
