// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/bxcodec/go-clean-arch/domain"
	mock "github.com/stretchr/testify/mock"
)

// TopicRepository is an autogenerated mock type for the TopicRepository type
type TopicRepository struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *TopicRepository) GetByID(ctx context.Context, id int64) (domain.Topic, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 domain.Topic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (domain.Topic, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.Topic); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Topic)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTopicRepository creates a new instance of TopicRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTopicRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TopicRepository {
	mock := &TopicRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
