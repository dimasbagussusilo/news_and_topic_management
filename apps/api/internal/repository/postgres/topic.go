package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
)

type TopicRepository struct {
	Conn *sql.DB
}

// NewTopicRepository will create an object that represent the topic.Repository interface
func NewTopicRepository(conn *sql.DB) *TopicRepository {
	return &TopicRepository{conn}
}

func (tr *TopicRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Topic, err error) {
	rows, err := tr.Conn.QueryContext(ctx, query, args...)
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

func (tr *TopicRepository) Fetch(ctx context.Context, filter domain.TopicFilter) (res []domain.Topic, totalData int64, err error) {
	query := `SELECT id, name, updated_at, created_at
			  FROM topic WHERE 1=1`
	countQuery := "SELECT COUNT(*) FROM topic WHERE 1=1"

	var args []interface{}
	argIndex := 1 // Start index for query parameters

	// Add conditions based on optional filters
	if filter.ID != 0 {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, filter.ID)
		argIndex++
	}
	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		countQuery += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, fmt.Sprintf("%%%s%%", filter.Name)) // Add wildcards
		argIndex++
	}

	// Execute the count query
	err = tr.Conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalData)
	if err != nil {
		return nil, 0, err
	}

	// Sorting logic
	if filter.SortBy != "" {
		orderDirection := "ASC" // default to ascending
		if filter.SortOrder == "desc" {
			orderDirection = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, orderDirection)
	}

	// Pagination logic
	if filter.Page < 1 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, offset)

	// Execute the main query with pagination
	res, err = tr.fetch(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return res, totalData, nil
}

func (tr *TopicRepository) GetByID(ctx context.Context, id int64) (res domain.Topic, err error) {
	query := `SELECT id, name, updated_at, created_at
			  FROM topic WHERE id = $1`

	list, err := tr.fetch(ctx, query, id)
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

func (tr *TopicRepository) GetByNewsID(ctx context.Context, id int64) (res domain.Topic, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
			  FROM topic WHERE id = $1`

	list, err := tr.fetch(ctx, query, id)
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

func (tr *TopicRepository) GetByName(ctx context.Context, name string) (res domain.Topic, err error) {
	query := `SELECT id, name, updated_at, created_at
			  FROM topic WHERE name = $1`

	list, err := tr.fetch(ctx, query, name)
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

func (tr *TopicRepository) Store(ctx context.Context, a *domain.Topic) (err error) {
	query := `INSERT INTO topic (name, updated_at, created_at)
			  VALUES ($1, $2, $3) RETURNING id`
	err = tr.Conn.QueryRowContext(ctx, query, a.Name, time.Now(), time.Now()).Scan(&a.ID)
	return
}

func (tr *TopicRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM topic WHERE id = $1"

	stmt, err := tr.Conn.PrepareContext(ctx, query)
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

func (tr *TopicRepository) Update(ctx context.Context, to *domain.Topic) (err error) {
	query := `UPDATE topic SET name=$1, updated_at=$2 WHERE id = $3`

	stmt, err := tr.Conn.PrepareContext(ctx, query)
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
