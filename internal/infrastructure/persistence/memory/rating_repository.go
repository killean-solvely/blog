package memory

import (
	"errors"
	"sync"

	"blog/internal/domain"
)

type RatingRepository struct {
	mu      sync.RWMutex
	ratings map[domain.RatingID]domain.Rating
}

func NewRatingRepository() *RatingRepository {
	return &RatingRepository{
		ratings: map[domain.RatingID]domain.Rating{},
	}
}

func (r *RatingRepository) All() ([]domain.Rating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ratings := []domain.Rating{}
	for k := range r.ratings {
		ratings = append(ratings, r.ratings[k])
	}

	return ratings, nil
}

func (r *RatingRepository) FindByID(id domain.RatingID) (*domain.Rating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rating, exists := r.ratings[id]
	if !exists {
		return nil, errors.New("no rows")
	}

	return &rating, nil
}

func (r *RatingRepository) FindByUser(userID domain.UserID) ([]domain.Rating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ratings := []domain.Rating{}
	for k := range r.ratings {
		if r.ratings[k].UserID() == userID {
			ratings = append(ratings, r.ratings[k])
		}
	}

	return ratings, nil
}

func (r *RatingRepository) FindByPost(postID domain.PostID) ([]domain.Rating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ratings := []domain.Rating{}
	for k := range r.ratings {
		if r.ratings[k].PostID() == postID {
			ratings = append(ratings, r.ratings[k])
		}
	}

	return ratings, nil
}

func (r *RatingRepository) Exists(id domain.RatingID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.ratings[id]
	return exists, nil
}

func (r *RatingRepository) Create(rating *domain.Rating) (*domain.Rating, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ratings[rating.GetID()] = *rating

	c := r.ratings[rating.GetID()]
	return &c, nil
}

func (r *RatingRepository) ChangeRating(id domain.RatingID, newRatingType domain.RatingType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c := r.ratings[id]
	c.ChangeRating(newRatingType)
	r.ratings[id] = c

	return nil
}

func (r *RatingRepository) RemoveRating(id domain.RatingID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.ratings[id]; !exists {
		return errors.New("doesn't exist")
	}

	delete(r.ratings, id)

	return nil
}
