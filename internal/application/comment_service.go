package application

import (
	"errors"
	"log"

	"blog/internal/domain"
	"blog/pkg/ddd"
)

type CommentService struct {
	commentRepo     domain.CommentRepository
	userRepo        domain.UserRepository
	postRepo        domain.PostRepository
	eventDispatcher ddd.EventDispatcher
}

func NewCommentService(
	commentRepo domain.CommentRepository,
	userRepo domain.UserRepository,
	postRepo domain.PostRepository,
	eventDispatcher ddd.EventDispatcher,
) *CommentService {
	return &CommentService{
		commentRepo:     commentRepo,
		userRepo:        userRepo,
		postRepo:        postRepo,
		eventDispatcher: eventDispatcher,
	}
}

func (s *CommentService) CreateComment(
	postID string,
	commenterID string,
	content string,
) (*CommentDTO, error) {
	domainPostID := domain.NewPostID(postID)
	domainCommenterID := domain.NewUserID(commenterID)

	// Check that the post exists
	if exists, err := s.postRepo.Exists(domainPostID); !exists || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("post does not exist")
	}

	// Check that the user exists
	if exists, err := s.userRepo.Exists(domainCommenterID); !exists || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("user does not exist")
	}

	// Create the comment
	comment, err := domain.NewComment(domainPostID, domainCommenterID, content)
	if err != nil {
		return nil, err
	}

	// Persist
	if _, err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(comment); err != nil {
		return nil, err
	}

	commentDTO := CommentDTO{}
	commentDTO.FromDomain(comment)

	return &commentDTO, nil
}

func (s *CommentService) EditComment(
	commentID string,
	content string,
) error {
	domainCommentID := domain.NewCommentID(commentID)

	// Check that the comment exists
	if exists, err := s.commentRepo.Exists(domainCommentID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("comment does not exist")
	}

	// Get and update the comment
	comment, err := s.commentRepo.FindByID(domainCommentID)
	if err != nil {
		return err
	}

	if err := comment.Edit(content); err != nil {
		return err
	}

	// Persist
	if err := s.commentRepo.UpdateContent(domainCommentID, content); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(comment); err != nil {
		return err
	}

	return nil
}

func (s *CommentService) ArchiveComment(
	commentID string,
) error {
	domainCommentID := domain.NewCommentID(commentID)

	// Check that the comment exists
	if exists, err := s.commentRepo.Exists(domainCommentID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("comment does not exist")
	}

	// Get and update the comment
	comment, err := s.commentRepo.FindByID(domainCommentID)
	if err != nil {
		return err
	}
	comment.Archive()

	// Persist
	if err := s.commentRepo.Archive(domainCommentID); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(comment); err != nil {
		return err
	}

	return nil
}

// Helper method to dispatch events for any aggregate with AggregateBase
func (s *CommentService) dispatchAggregateEvents(aggregate ddd.EventAggregate) error {
	events := aggregate.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventDispatcher.Dispatch(event); err != nil {
			log.Printf("Failed to dispatch event: %v", err)
		}
	}
	aggregate.MarkEventsAsCommitted()
	return nil
}
