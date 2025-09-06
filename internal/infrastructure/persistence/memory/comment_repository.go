package memory

import (
	"errors"
	"sync"

	"blog/internal/domain"
)

type CommentRepository struct {
	mu       sync.RWMutex
	comments map[domain.CommentID]domain.Comment
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		comments: map[domain.CommentID]domain.Comment{},
	}
}

func (r *CommentRepository) All() ([]domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := []domain.Comment{}
	for k := range r.comments {
		comments = append(comments, r.comments[k])
	}

	return comments, nil
}

func (r *CommentRepository) FindByID(id domain.CommentID) (*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, errors.New("no rows")
	}

	return &comment, nil
}

func (r *CommentRepository) FindByUser(userID domain.UserID) ([]domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := []domain.Comment{}
	for k := range r.comments {
		if r.comments[k].CommenterID() == userID {
			comments = append(comments, r.comments[k])
		}
	}

	return comments, nil
}

func (r *CommentRepository) FindByPost(postID domain.PostID) ([]domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comments := []domain.Comment{}
	for k := range r.comments {
		if r.comments[k].PostID() == postID {
			comments = append(comments, r.comments[k])
		}
	}

	return comments, nil
}

func (r *CommentRepository) Exists(id domain.CommentID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.comments[id]
	return exists, nil
}

func (r *CommentRepository) Create(comment *domain.Comment) (*domain.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.comments[comment.GetID()] = *comment

	c := r.comments[comment.GetID()]
	return &c, nil
}

func (r *CommentRepository) UpdateContent(id domain.CommentID, newContent string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c := r.comments[id]
	c.Edit(newContent)
	r.comments[id] = c

	return nil
}

func (r *CommentRepository) Archive(id domain.CommentID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c := r.comments[id]
	c.Archive()
	r.comments[id] = c

	return nil
}
