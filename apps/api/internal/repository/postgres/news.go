package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/repository"
)

type NewsRepository struct {
	Conn *sql.DB
}

// NewNewsRepository will create an object that represent the news.Repository interface
func NewNewsRepository(conn *sql.DB) *NewsRepository {
	return &NewsRepository{conn}
}

func (m *NewsRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.News, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(rows)

	result = make([]domain.News, 0)
	for rows.Next() {
		t := domain.News{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&t.Author.ID,
			&t.Status,
			&t.UpdatedAt,
			&t.CreatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (m *NewsRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.News, nextCursor string, err error) {
	query := `SELECT id, title, content, author_id, status, updated_at, created_at
			  FROM news WHERE created_at > $1 ORDER BY created_at LIMIT $2`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}
	return
}

func (m *NewsRepository) GetByID(ctx context.Context, id int64) (res domain.News, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
			  FROM news WHERE id = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.News{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *NewsRepository) GetByTitle(ctx context.Context, title string) (res domain.News, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
			  FROM news WHERE title = $1`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *NewsRepository) Store(ctx context.Context, a *domain.News) (err error) {
	query := `INSERT INTO news (title, content, author_id, updated_at, created_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt).Scan(&a.ID)
	return
}

func (m *NewsRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM news WHERE id = $1"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("unexpected behavior: total affected rows = %d", rowsAffected)
		return
	}

	return
}

func (m *NewsRepository) Update(ctx context.Context, ar *domain.News) (err error) {
	query := `UPDATE news SET title=$1, content=$2, author_id=$3, updated_at=$4 WHERE id = $5`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affected != 1 {
		err = fmt.Errorf("unexpected behavior: total affected rows = %d", affected)
		return
	}

	return
}
