package topic

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
)

// TopicRepository represent the news's repository contract
//
//go:generate mockery --name TopicRepository
type TopicRepository interface {
	Fetch(ctx context.Context, filter domain.TopicFilter) (res []domain.Topic, totalPage int64, err error)
	GetByName(ctx context.Context, name string) (domain.Topic, error)
	GetByID(ctx context.Context, id int64) (domain.Topic, error)
	Update(ctx context.Context, ar *domain.Topic) error
	Store(ctx context.Context, a *domain.Topic) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	topicRepo TopicRepository
}

// NewService will create a new topic service object
func NewService(t TopicRepository) *Service {
	return &Service{
		topicRepo: t,
	}
}

/*
* In this function below, I'm using err group with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */

func (s *Service) Fetch(ctx context.Context, filter domain.TopicFilter) (res []domain.Topic, totalPage int64, err error) {
	res, totalPage, err = s.topicRepo.Fetch(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return
}

func (s *Service) GetByID(ctx context.Context, id int64) (res domain.Topic, err error) {
	res, err = s.topicRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (s *Service) Update(ctx context.Context, unr *domain.Topic) (err error) {
	return s.topicRepo.Update(ctx, unr)
}

func (s *Service) GetByTitle(ctx context.Context, name string) (res domain.Topic, err error) {
	res, err = s.topicRepo.GetByName(ctx, name)
	if err != nil {
		return
	}

	return
}

func (s *Service) Store(ctx context.Context, cnr *domain.Topic) (err error) {
	existedTopic, _ := s.GetByTitle(ctx, cnr.Name) // ignore if any error
	if existedTopic.ID != 0 {
		return domain.ErrConflict
	}

	err = s.topicRepo.Store(ctx, cnr)
	if err != nil {
		return
	}

	return
}

func (s *Service) Delete(ctx context.Context, id int64) (err error) {
	existedTopic, err := s.topicRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	if existedTopic.ID == 0 {
		return domain.ErrNotFound
	}
	return s.topicRepo.Delete(ctx, id)
}
