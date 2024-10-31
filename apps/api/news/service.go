package news

import (
	"context"
	"fmt"
	"github.com/bxcodec/go-clean-arch/internal/dto/news"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/bxcodec/go-clean-arch/domain"
)

// NewsRepository represent the news's repository contract
//
//go:generate mockery --name NewsRepository
type NewsRepository interface {
	Fetch(ctx context.Context, filter domain.NewsFilter) (res []domain.News, totalPage int64, err error)
	GetByID(ctx context.Context, id int64) (domain.News, error)
	GetByTitle(ctx context.Context, title string) (domain.News, error)
	Update(ctx context.Context, ar *news.UpdateNewsReq) error
	Store(ctx context.Context, a *news.CreateNewsReq) error
	Delete(ctx context.Context, id int64) error
}

// AuthorRepository represent the author's repository contract
//
//go:generate mockery --name AuthorRepository
type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Author, error)
}

// TopicRepository represent topics repository contract
//
//go:generate mockery --name TopicRepository
type TopicRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Topic, error)
}

// NewsTopicRepository represent the news topic repository contract
//
//go:generate mockery --name NewsTopicRepository
type NewsTopicRepository interface {
	GetByNewsID(ctx context.Context, newsId int64) ([]domain.NewsTopic, error)
	GetByTopicID(ctx context.Context, topicId int64) ([]domain.NewsTopic, error)
	Store(ctx context.Context, nt *domain.NewsTopic) (err error)
	DeleteByNewsID(ctx context.Context, newsId int64) (err error)
}

type Service struct {
	newsRepo      NewsRepository
	authorRepo    AuthorRepository
	topicRepo     TopicRepository
	newsTopicRepo NewsTopicRepository
}

// NewService will create a new news service object
func NewService(n NewsRepository, a AuthorRepository, t TopicRepository, nt NewsTopicRepository) *Service {
	return &Service{
		newsRepo:      n,
		authorRepo:    a,
		topicRepo:     t,
		newsTopicRepo: nt,
	}
}

/*
* In this function below, I'm using err group with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (s *Service) fillAuthorDetails(ctx context.Context, data []domain.News) ([]domain.News, error) {
	g, ctx := errgroup.WithContext(ctx)
	// Get the author's id
	mapAuthors := map[int64]domain.Author{}

	for _, newsData := range data { //nolint
		mapAuthors[newsData.Author.ID] = domain.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan domain.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := s.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	go func() {
		defer close(chanAuthor)
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}

	}()

	for author := range chanAuthor {
		if author != (domain.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data { //nolint
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = domain.AuthorNews{
				ID:   a.ID,
				Name: a.Name,
			}
		}
	}
	return data, nil
}

func (s *Service) fillTopicDetails(ctx context.Context, data []domain.News) ([]domain.News, error) {
	g, ctx := errgroup.WithContext(ctx)
	mapNewsTopics := map[int64][]domain.Topic{}

	chanNewsTopic := make(chan []domain.NewsTopic)
	for _, newsData := range data {
		newsID := newsData.ID
		g.Go(func() error {
			res, err := s.newsTopicRepo.GetByNewsID(ctx, newsID)
			if err != nil {
				return err
			}
			chanNewsTopic <- res
			return nil
		})
	}

	go func() {
		defer close(chanNewsTopic)
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
	}()

	for newsTopic := range chanNewsTopic {
		for _, nt := range newsTopic {
			mapNewsTopics[nt.NewsID] = append(mapNewsTopics[nt.NewsID], domain.Topic{ID: nt.TopicID})
		}
	}

	chanTopic := make(chan domain.Topic)
	for _, topics := range mapNewsTopics {
		for _, topic := range topics {
			topicID := topic.ID
			g.Go(func() error {
				res, err := s.topicRepo.GetByID(context.Background(), topicID)
				if err != nil {
					return err
				}
				chanTopic <- res
				return nil
			})
		}
	}

	go func() {
		defer close(chanTopic)
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
	}()

	for topic := range chanTopic {
		for newsID, topicList := range mapNewsTopics {
			for i, t := range topicList {
				if t.ID == topic.ID {
					mapNewsTopics[newsID][i] = topic
				}
			}
		}
	}

	for i, newsData := range data {
		if topics, ok := mapNewsTopics[newsData.ID]; ok {
			var topicNews []domain.TopicNews
			for _, topic := range topics {
				topicNews = append(topicNews, domain.TopicNews{
					ID:   topic.ID,
					Name: topic.Name,
				})
			}
			data[i].Topics = topicNews
		}
	}

	return data, nil
}

func (s *Service) Fetch(ctx context.Context, filter domain.NewsFilter) (res []domain.News, totalPage int64, err error) {
	res, totalPage, err = s.newsRepo.Fetch(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if res, err = s.fillAuthorDetails(ctx, res); err != nil {
		totalPage = 0
		return
	}

	if res, err = s.fillTopicDetails(ctx, res); err != nil {
		totalPage = 0
		return
	}
	return
}

func (s *Service) GetByID(ctx context.Context, id int64) (res domain.News, err error) {
	res, err = s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := s.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.News{}, err
	}

	res.Author = domain.AuthorNews{
		ID:   resAuthor.ID,
		Name: resAuthor.Name,
	}
	return
}

func (s *Service) Update(ctx context.Context, unr *news.UpdateNewsReq) (err error) {
	if unr.TopicIDs != nil {
		// Remove previous news topics
		_ = s.newsTopicRepo.DeleteByNewsID(ctx, *unr.ID)

		// Create new news topics
		for _, topicId := range *unr.TopicIDs {
			err = s.newsTopicRepo.Store(ctx, &domain.NewsTopic{
				NewsID:  *unr.ID,
				TopicID: topicId,
			})
			if err != nil {
				return fmt.Errorf("failed to store new news topic: %w", err)
			}
		}
	}

	// Step 3: Update the news article itself
	return s.newsRepo.Update(ctx, unr)
}

func (s *Service) GetByTitle(ctx context.Context, title string) (res domain.News, err error) {
	res, err = s.newsRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := s.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.News{}, err
	}

	res.Author = domain.AuthorNews{
		ID:   resAuthor.ID,
		Name: resAuthor.Name,
	}
	return
}

func (s *Service) Store(ctx context.Context, cnr *news.CreateNewsReq) (err error) {
	existedNews, _ := s.GetByTitle(ctx, cnr.Title) // ignore if any error
	if existedNews.ID != 0 {
		return domain.ErrConflict
	}

	err = s.newsRepo.Store(ctx, cnr)
	if err != nil {
		return
	}
	for _, topicId := range cnr.TopicIDs {
		err = s.newsTopicRepo.Store(ctx, &domain.NewsTopic{
			NewsID:  cnr.ID,
			TopicID: topicId,
		})
		if err != nil {
			return
		}

	}
	return
}

func (s *Service) Delete(ctx context.Context, id int64) (err error) {
	existedNews, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	if existedNews.ID == 0 {
		return domain.ErrNotFound
	}
	return s.newsRepo.Delete(ctx, id)
}
