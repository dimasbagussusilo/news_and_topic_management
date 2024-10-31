package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/internal/repository"
)

type TopicRepository struct {
	Conn *sql.DB
}

// NewTopicRepository will create an object that represent the topic.Repository interface
func NewTopicRepository(conn *sql.DB) *TopicRepository {
	return &TopicRepository{conn}
}

func (m *TopicRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Topic, err error) {
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

	result = make([]domain.Topic, 0)
	for rows.Next() {
		t := domain.Topic{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
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

func (m *TopicRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Topic, nextCursor string, err error) {
	query := `SELECT id, name, updated_at, created_at
			  FROM topic WHERE created_at > $1 ORDER BY created_at LIMIT $2`

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

func (m *TopicRepository) GetByID(ctx context.Context, id int64) (res domain.Topic, err error) {
	query := `SELECT id, name, updated_at, created_at
			  FROM topic WHERE id = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Topic{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *TopicRepository) GetByNewsID(ctx context.Context, id int64) (res domain.Topic, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
			  FROM topic WHERE id = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Topic{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *TopicRepository) Store(ctx context.Context, a *domain.Topic) (err error) {
	query := `INSERT INTO topic (name, updated_at, created_at)
			  VALUES ($1, $2, $3) RETURNING id`
	err = m.Conn.QueryRowContext(ctx, query, a.Name, a.UpdatedAt, a.CreatedAt).Scan(&a.ID)
	return
}

func (m *TopicRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM topic WHERE id = $1"

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

func (m *TopicRepository) Update(ctx context.Context, to *domain.Topic) (err error) {
	query := `UPDATE topic SET name=$1, updated_at=$2 WHERE id = $3`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, to.Name, to.UpdatedAt, to.ID)
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
