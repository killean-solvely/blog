package application

import (
	"errors"
	"log"

	"blog/internal/domain"
	"blog/pkg/ddd"
)

type PostService struct {
	postRepo        domain.PostRepository
	userRepo        domain.UserRepository
	eventDispatcher ddd.EventDispatcher
}

func NewPostService(
	postRepo domain.PostRepository,
	userRepo domain.UserRepository,
	dispatcher ddd.EventDispatcher,
) *PostService {
	return &PostService{
		postRepo:        postRepo,
		userRepo:        userRepo,
		eventDispatcher: dispatcher,
	}
}

func (s *PostService) CreatePost(
	authorID string,
	title string,
	content string,
) (*PostDTO, error) {
	domainAuthorID := domain.NewUserID(authorID)

	// Check that the user exists
	if exists, err := s.userRepo.Exists(domainAuthorID); !exists || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("author does not exist")
	}

	// Create the post
	post, err := domain.NewPost(domainAuthorID, title, content)
	if err != nil {
		return nil, err
	}

	// Persist
	if _, err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(post); err != nil {
		return nil, err
	}

	postDTO := PostDTO{}
	postDTO.FromDomain(post)

	return &postDTO, nil
}

// Helper method to dispatch events for any aggregate with AggregateBase
func (s *PostService) dispatchAggregateEvents(aggregate ddd.EventAggregate) error {
	events := aggregate.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventDispatcher.Dispatch(event); err != nil {
			log.Printf("Failed to dispatch event: %v", err)
		}
	}
	aggregate.MarkEventsAsCommitted()
	return nil
}
