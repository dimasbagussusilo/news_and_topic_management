package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/repository"
	newsPostgresRepo "github.com/bxcodec/go-clean-arch/internal/repository/postgres"
)

func TestFetchNews(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockNews := []domain.News{
		{
			ID: 1, Title: "title 1", Content: "content 1",
			Author: domain.AuthorNews{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		{
			ID: 2, Title: "title 2", Content: "content 2",
			Author: domain.AuthorNews{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(mockNews[0].ID, mockNews[0].Title, mockNews[0].Content,
			mockNews[0].Author.ID, mockNews[0].UpdatedAt, mockNews[0].CreatedAt).
		AddRow(mockNews[1].ID, mockNews[1].Title, mockNews[1].Content,
			mockNews[1].Author.ID, mockNews[1].UpdatedAt, mockNews[1].CreatedAt)

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM news WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := newsPostgresRepo.NewNewsRepository(db)
	cursor := repository.EncodeCursor(mockNews[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetNewsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM news WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := newsPostgresRepo.NewNewsRepository(db)

	num := int64(5)
	anNews, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anNews)
}

func TestStoreNews(t *testing.T) {
	now := time.Now()
	ar := &domain.News{
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: domain.AuthorNews{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT  news SET title=\\? , content=\\? , author_id=\\?, updated_at=\\? , created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.CreatedAt, ar.UpdatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := newsPostgresRepo.NewNewsRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}

func TestGetNewsByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, author_id, updated_at, created_at FROM news WHERE title = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := newsPostgresRepo.NewNewsRepository(db)

	title := "title 1"
	anNews, err := a.GetByTitle(context.TODO(), title)
	assert.NoError(t, err)
	assert.NotNil(t, anNews)
}

func TestDeleteNews(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM news WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	a := newsPostgresRepo.NewNewsRepository(db)

	num := int64(12)
	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}

func TestUpdateNews(t *testing.T) {
	now := time.Now()
	ar := &domain.News{
		ID:        12,
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: domain.AuthorNews{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE news set title=\\?, content=\\?, author_id=\\?, updated_at=\\? WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	a := newsPostgresRepo.NewNewsRepository(db)

	err = a.Update(context.TODO(), ar)
	assert.NoError(t, err)
}
