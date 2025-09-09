package application

import (
	"errors"
	"log"

	"blog/internal/domain"
	"blog/pkg/ddd"
)

type RatingService struct {
	ratingRepo      domain.RatingRepository
	userRepo        domain.UserRepository
	postRepo        domain.PostRepository
	eventDispatcher ddd.EventDispatcher
}

func NewRatingService(
	ratingRepo domain.RatingRepository,
	userRepo domain.UserRepository,
	postRepo domain.PostRepository,
	eventDispatcher ddd.EventDispatcher,
) *RatingService {
	return &RatingService{
		ratingRepo:      ratingRepo,
		userRepo:        userRepo,
		postRepo:        postRepo,
		eventDispatcher: eventDispatcher,
	}
}

func (s RatingService) GetRatingsOnPost(postID string) ([]RatingDTO, error) {
	domainPostID := domain.NewPostID(postID)

	ratings, err := s.ratingRepo.FindByPost(domainPostID)
	if err != nil {
		return nil, err
	}

	ratingDTOs := []RatingDTO{}
	for i := range ratings {
		ratingDTO := RatingDTO{}
		ratingDTO.FromDomain(&ratings[i])
		ratingDTOs = append(ratingDTOs, ratingDTO)
	}

	return ratingDTOs, nil
}

func (s RatingService) GetRating(ratingID string) (*RatingDTO, error) {
	domainRatingID := domain.NewRatingID(ratingID)

	rating, err := s.ratingRepo.FindByID(domainRatingID)
	if err != nil {
		return nil, err
	}

	ratingDTO := RatingDTO{}
	ratingDTO.FromDomain(rating)

	return &ratingDTO, nil
}

func (s *RatingService) CreateRating(
	postID string,
	userID string,
	ratingType string,
) (*RatingDTO, error) {
	domainPostID := domain.NewPostID(postID)
	domainUserID := domain.NewUserID(userID)

	// Check that the post exists
	if exists, err := s.postRepo.Exists(domainPostID); !exists || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("post does not exist")
	}

	// Check that the user exists
	if exists, err := s.userRepo.Exists(domainUserID); !exists || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("user does not exist")
	}

	// Check if rating already exists for this user/post combination
	if exists, err := s.ratingRepo.ExistsOnPostByUser(domainPostID, domainUserID); !exists ||
		err != nil {
		return nil, errors.New("rating already exists for this user and post")
	}

	// Create the rating
	rating := domain.NewRating(domainPostID, domainUserID, domain.RatingType(ratingType))

	// Persist
	if _, err := s.ratingRepo.Create(rating); err != nil {
		return nil, err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(rating); err != nil {
		return nil, err
	}

	ratingDTO := RatingDTO{}
	ratingDTO.FromDomain(rating)

	return &ratingDTO, nil
}

func (s *RatingService) UpdateRating(
	ratingID string,
	newRatingType string,
) error {
	domainRatingID := domain.NewRatingID(ratingID)
	domainRatingType := domain.RatingType(newRatingType)

	// Check that the rating exists
	if exists, err := s.ratingRepo.Exists(domainRatingID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("rating does not exist")
	}

	// Get the rating and update it
	rating, err := s.ratingRepo.FindByID(domainRatingID)
	if err != nil {
		return err
	}

	rating.ChangeRating(domainRatingType)

	// Persist
	if err := s.ratingRepo.ChangeRating(domainRatingID, domainRatingType); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(rating); err != nil {
		return err
	}

	return nil
}

func (s *RatingService) RemoveRating(
	ratingID string,
) error {
	domainRatingID := domain.NewRatingID(ratingID)

	// Check that the rating exists
	if exists, err := s.ratingRepo.Exists(domainRatingID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("rating does not exist")
	}

	// Get the rating and remove it
	rating, err := s.ratingRepo.FindByID(domainRatingID)
	if err != nil {
		return err
	}

	rating.RemoveRating()

	// Persist (delete the rating)
	if err := s.ratingRepo.RemoveRating(domainRatingID); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(rating); err != nil {
		return err
	}

	return nil
}

// Helper method to dispatch events for any aggregate with AggregateBase
func (s *RatingService) dispatchAggregateEvents(aggregate ddd.EventAggregate) error {
	events := aggregate.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventDispatcher.Dispatch(event); err != nil {
			log.Printf("Failed to dispatch event: %v", err)
		}
	}
	aggregate.MarkEventsAsCommitted()
	return nil
}
