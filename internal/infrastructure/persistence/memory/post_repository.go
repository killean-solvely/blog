package memory

import (
	"errors"
	"sync"

	"blog/internal/domain"
)

type PostRepository struct {
	mu    sync.RWMutex
	posts map[domain.PostID]domain.Post
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts: map[domain.PostID]domain.Post{},
	}
}

func (r *PostRepository) All() ([]domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := []domain.Post{}
	for k := range r.posts {
		posts = append(posts, r.posts[k])
	}

	return posts, nil
}

func (r *PostRepository) FindByID(id domain.PostID) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("no rows")
	}

	return &post, nil
}

func (r *PostRepository) FindByAuthor(userID domain.UserID) ([]domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := []domain.Post{}
	for k := range r.posts {
		if r.posts[k].AuthorID() == userID {
			posts = append(posts, r.posts[k])
		}
	}

	return posts, nil
}

func (r *PostRepository) Create(post *domain.Post) (*domain.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.GetID()] = *post

	p := r.posts[post.GetID()]
	return &p, nil
}

func (r *PostRepository) UpdateTitle(id domain.PostID, newTitle string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.posts[id]
	p.EditTitle(newTitle)
	r.posts[id] = p

	return nil
}

func (r *PostRepository) UpdateContent(id domain.PostID, newContent string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.posts[id]
	p.EditContent(newContent)
	r.posts[id] = p

	return nil
}

func (r *PostRepository) Archive(id domain.PostID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := r.posts[id]
	p.Archive()
	r.posts[id] = p

	return nil
}
