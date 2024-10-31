package news_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/news"
	"github.com/bxcodec/go-clean-arch/news/mocks"
)

func TestFetch(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	mockTopicRepo := new(mocks.TopicRepository)
	mockNewsTopicRepo := new(mocks.NewsTopicRepository)
	mockNews := domain.News{
		Title:   "Hello",
		Content: "Content",
	}

	mockListArticle := []domain.News{mockNews}

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArticle, "next-cursor", nil).Once()

		mockAuthor := domain.Author{ID: 1, Name: "Iman Tumorang"}
		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)
		cursor := "12"
		num := int64(1)

		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.NoError(t, err)
		assert.Equal(t, "next-cursor", nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.Len(t, list, len(mockListArticle))

		mockNewsRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockNewsRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpected Error")).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)
		cursor := "12"
		num := int64(1)

		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Error(t, err)
		assert.Empty(t, nextCursor)
		assert.Len(t, list, 0)

		mockNewsRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	mockTopicRepo := new(mocks.TopicRepository)
	mockNewsTopicRepo := new(mocks.NewsTopicRepository)
	mockNews := domain.News{Title: "Hello", Content: "Content"}
	mockAuthor := domain.Author{ID: 1, Name: "Iman Tumorang"}

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockNews, nil).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)
		a, err := u.GetByID(context.TODO(), mockNews.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockNewsRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockNewsRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.News{}, errors.New("Unexpected")).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)

		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)
		a, err := u.GetByID(context.TODO(), mockNews.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.News{}, a)

		mockNewsRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestStore(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	mockTopicRepo := new(mocks.TopicRepository)
	mockNewsTopicRepo := new(mocks.NewsTopicRepository)
	mockNews := domain.News{Title: "Hello", Content: "Content"}

	t.Run("success", func(t *testing.T) {
		tempMockNews := mockNews
		tempMockNews.ID = 0
		mockNewsRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(domain.News{}, domain.ErrNotFound).Once()
		mockNewsRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.News")).Return(nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)

		err := u.Store(context.TODO(), &tempMockNews)

		assert.NoError(t, err)
		assert.Equal(t, mockNews.Title, tempMockNews.Title)
		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("existing-title", func(t *testing.T) {
		existingNews := mockNews
		mockNewsRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingNews, nil).Once()

		mockAuthor := domain.Author{ID: 1, Name: "Iman Tumorang"}
		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)
		err := u.Store(context.TODO(), &mockNews)

		assert.Error(t, err)
		mockNewsRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	mockTopicRepo := new(mocks.TopicRepository)
	mockNewsTopicRepo := new(mocks.NewsTopicRepository)
	mockNews := domain.News{Title: "Hello", Content: "Content"}

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockNews, nil).Once()
		mockNewsRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)

		err := u.Delete(context.TODO(), mockNews.ID)

		assert.NoError(t, err)
		mockNewsRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("news-is-not-exist", func(t *testing.T) {
		mockNewsRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.News{}, nil).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)
		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)

		err := u.Delete(context.TODO(), mockNews.ID)

		assert.Error(t, err)
		mockNewsRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error-happens-in-db", func(t *testing.T) {
		mockNewsRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.News{}, errors.New("Unexpected Error")).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)

		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)
		err := u.Delete(context.TODO(), mockNews.ID)

		assert.Error(t, err)
		mockNewsRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	mockTopicRepo := new(mocks.TopicRepository)
	mockNewsTopicRepo := new(mocks.NewsTopicRepository)
	mockNews := domain.News{Title: "Hello", Content: "Content", ID: 23}

	t.Run("success", func(t *testing.T) {
		mockNewsRepo.On("Update", mock.Anything, &mockNews).Return(nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := news.NewService(mockNewsRepo, mockAuthorRepo, mockTopicRepo, mockNewsTopicRepo)

		err := u.Update(context.TODO(), &mockNews)

		assert.NoError(t, err)
		mockNewsRepo.AssertExpectations(t)
	})
}
