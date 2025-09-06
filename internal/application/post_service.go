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

func (s *PostService) UpdatePostTitle(postID string, newTitle string) error {
	domainPostID := domain.NewPostID(postID)

	// Check that the post exists
	if exists, err := s.postRepo.Exists(domainPostID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("post does not exist")
	}

	// Get the post, then update it
	post, err := s.postRepo.FindByID(domainPostID)
	if err != nil {
		return err
	}

	if err := post.EditTitle(newTitle); err != nil {
		return err
	}

	// Persist
	if err := s.postRepo.UpdateTitle(domainPostID, newTitle); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(post); err != nil {
		return err
	}

	return nil
}

func (s *PostService) UpdatePostContent(postID string, newContent string) error {
	domainPostID := domain.NewPostID(postID)

	// Check that the post exists
	if exists, err := s.postRepo.Exists(domainPostID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("post does not exist")
	}

	// Get the post, then update it
	post, err := s.postRepo.FindByID(domainPostID)
	if err != nil {
		return err
	}

	if err := post.EditContent(newContent); err != nil {
		return err
	}

	// Persist
	if err := s.postRepo.UpdateContent(domainPostID, newContent); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(post); err != nil {
		return err
	}

	return nil
}

func (s *PostService) ArchivePost(postID string) error {
	domainPostID := domain.NewPostID(postID)

	// Make sure the post exists first
	if exists, err := s.postRepo.Exists(domainPostID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("post does not exist")
	}

	// Get the post, then update it
	post, err := s.postRepo.FindByID(domainPostID)
	if err != nil {
		return err
	}

	post.Archive()

	// Persist
	if err := s.postRepo.Archive(domainPostID); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(post); err != nil {
		return err
	}

	return nil
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
