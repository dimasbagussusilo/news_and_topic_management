package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bxcodec/go-clean-arch/internal/dto/news"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/sirupsen/logrus"
)

type NewsRepository struct {
	Conn *sql.DB
}

// NewNewsRepository will create an object that represent the news.Repository interface
func NewNewsRepository(conn *sql.DB) *NewsRepository {
	return &NewsRepository{conn}
}

func (nr *NewsRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.News, err error) {
	rows, err := nr.Conn.QueryContext(ctx, query, args...)
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

func (nr *NewsRepository) Fetch(ctx context.Context, filter domain.NewsFilter) (res []domain.News, totalData int64, err error) {
	query := `SELECT id, title, content, author_id, status, updated_at, created_at
			  FROM news WHERE 1=1`
	countQuery := "SELECT COUNT(*) FROM news WHERE 1=1"

	var args []interface{}
	argIndex := 1 // Start index for query parameters

	// Add conditions based on optional filters
	if filter.ID != 0 {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, filter.ID)
		argIndex++
	}
	if filter.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		countQuery += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		args = append(args, fmt.Sprintf("%%%s%%", filter.Title)) // Add wildcards
		argIndex++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}
	if filter.AuthorID != 0 {
		query += fmt.Sprintf(" AND author_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND author_id = $%d", argIndex)
		args = append(args, filter.AuthorID)
		argIndex++
	}
	var zeroTime time.Time
	if filter.StartDate != zeroTime {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, filter.StartDate)
		argIndex++
	}
	if filter.EndDate != zeroTime {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, filter.EndDate)
		argIndex++
	}

	// Execute the count query
	err = nr.Conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalData)
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
	res, err = nr.fetch(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return res, totalData, nil
}

func (nr *NewsRepository) GetByID(ctx context.Context, id int64) (res domain.News, err error) {
	query := `SELECT id, title, content, author_id, status, updated_at, created_at
			  FROM news WHERE id = $1`

	list, err := nr.fetch(ctx, query, id)
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

func (nr *NewsRepository) GetByTitle(ctx context.Context, title string) (res domain.News, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
			  FROM news WHERE title = $1`

	list, err := nr.fetch(ctx, query, title)
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

func (nr *NewsRepository) Store(ctx context.Context, n *news.CreateNewsReq) (err error) {
	query := `INSERT INTO news (title, content, author_id, status, updated_at, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = nr.Conn.QueryRowContext(ctx, query, n.Title, n.Content, n.AuthorID, n.Status, time.Now(), time.Now()).Scan(&n.ID)
	return
}

func (nr *NewsRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM news WHERE id = $1"

	stmt, err := nr.Conn.PrepareContext(ctx, query)
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

func (nr *NewsRepository) Update(ctx context.Context, cnr *news.UpdateNewsReq) (err error) {
	// Ensure cnr.ID is provided
	if cnr.ID == nil {
		return errors.New("news ID is required for update")
	}

	// Initialize the base query and an args slice
	query := `UPDATE news SET `
	var args []interface{}
	argIndex := 1

	// Dynamically build the update query based on which fields are set
	if cnr.Title != nil {
		query += fmt.Sprintf("title = $%d, ", argIndex)
		args = append(args, *cnr.Title)
		argIndex++
	}
	if cnr.Content != nil {
		query += fmt.Sprintf("content = $%d, ", argIndex)
		args = append(args, *cnr.Content)
		argIndex++
	}
	if cnr.AuthorID != nil {
		query += fmt.Sprintf("author_id = $%d, ", argIndex)
		args = append(args, *cnr.AuthorID)
		argIndex++
	}
	if cnr.Status != nil {
		query += fmt.Sprintf("status = $%d, ", argIndex)
		args = append(args, *cnr.Status)
		argIndex++
	}

	// Always update the updated_at timestamp
	query += fmt.Sprintf("updated_at = $%d ", argIndex)
	args = append(args, time.Now())
	argIndex++

	// Complete the WHERE clause
	query += fmt.Sprintf("WHERE id = $%d", argIndex)
	args = append(args, *cnr.ID)

	// Prepare and execute the query
	stmt, err := nr.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			return
		}
	}(stmt) // Ensure statement is closed after execution

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected != 1 {
		err = fmt.Errorf("unexpected behavior: total affected rows = %d", affected)
		return
	}

	return nil
}
