package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/domain"
)

type NewsTopicRepository struct {
	Conn *sql.DB
}

// NewNewsTopicRepository will create an object that represent the newsTopic.Repository interface
func NewNewsTopicRepository(conn *sql.DB) *NewsTopicRepository {
	return &NewsTopicRepository{conn}
}

func (ntr *NewsTopicRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.NewsTopic, err error) {
	rows, err := ntr.Conn.QueryContext(ctx, query, args...)
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

	result = make([]domain.NewsTopic, 0)
	for rows.Next() {
		t := domain.NewsTopic{}
		err = rows.Scan(
			&t.NewsID,
			&t.TopicID,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (ntr *NewsTopicRepository) GetByNewsID(ctx context.Context, newsId int64) (res []domain.NewsTopic, err error) {
	query := `SELECT news_id, topic_id
			  FROM news_topic WHERE news_id = $1`

	list, err := ntr.fetch(ctx, query, newsId)
	if err != nil {
		return []domain.NewsTopic{}, err
	}

	if len(list) == 0 {
		return res, domain.ErrNotFound
	}

	return list, nil
}

func (ntr *NewsTopicRepository) GetByTopicID(ctx context.Context, topicId int64) (res []domain.NewsTopic, err error) {
	query := `SELECT news_id, topic_id
			  FROM news_topic WHERE topic_id = $1`

	list, err := ntr.fetch(ctx, query, topicId)
	if err != nil {
		return []domain.NewsTopic{}, err
	}

	if len(list) == 0 {
		return res, domain.ErrNotFound
	}

	return list, nil
}

func (ntr *NewsTopicRepository) Store(ctx context.Context, nt *domain.NewsTopic) (err error) {
	query := `INSERT INTO news_topic (news_id, topic_id)
			  VALUES ($1, $2) RETURNING news_id, topic_id`
	err = ntr.Conn.QueryRowContext(ctx, query, nt.NewsID, nt.TopicID).Scan(&nt.NewsID, &nt.TopicID)
	return
}

func (ntr *NewsTopicRepository) Delete(ctx context.Context, newsId int64, topicId int64) (err error) {
	query := "DELETE FROM news_topic WHERE news_id = $1 AND topic_id = $2"

	stmt, err := ntr.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, newsId, topicId)
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

func (ntr *NewsTopicRepository) DeleteByNewsID(ctx context.Context, newsId int64) (err error) {
	query := "DELETE FROM news_topic WHERE news_id = $1"

	stmt, err := ntr.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, newsId)
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
