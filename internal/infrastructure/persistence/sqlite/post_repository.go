package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"blog/internal/domain"
	"blog/internal/infrastructure/persistence/models"

	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r PostRepository) All() ([]domain.Post, error) {
	var dbPosts []models.Post
	err := r.db.Select(&dbPosts, "SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	posts := dbPostsToDomainPosts(dbPosts)
	return posts, nil
}

func (r PostRepository) FindByID(id domain.PostID) (*domain.Post, error) {
	var dbPost models.Post
	err := r.db.Get(&dbPost, "SELECT * FROM posts WHERE id=?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, err
	}

	post := dbPostToDomainPost(dbPost)
	return post, nil
}

func (r PostRepository) FindByAuthor(authorID domain.UserID) ([]domain.Post, error) {
	var dbPosts []models.Post
	err := r.db.Select(&dbPosts, "SELECT * FROM posts WHERE author_id=?", authorID)
	if err != nil {
		return nil, err
	}

	posts := dbPostsToDomainPosts(dbPosts)
	return posts, nil
}

func (r PostRepository) Exists(id domain.PostID) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM posts WHERE id=?", id)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r PostRepository) Create(post *domain.Post) (*domain.Post, error) {
	_, err := r.db.Exec(`
		INSERT INTO 
		posts (id, author_id, title, content, created_at, last_edited_at, archived_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		post.GetID().String(),
		post.AuthorID().String(),
		post.Title(),
		post.Content(),
		post.CreatedAt(),
		post.LastEditedAt(),
		post.ArchivedAt(),
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r PostRepository) UpdateTitle(id domain.PostID, newTitle string) error {
	_, err := r.db.Exec(`
		UPDATE posts
		SET title = ?, last_edited_at = ?
		WHERE id = ?
	`,
		newTitle,
		time.Now(),
		id.String(),
	)
	return err
}

func (r PostRepository) UpdateContent(id domain.PostID, newContent string) error {
	_, err := r.db.Exec(`
		UPDATE posts
		SET content = ?, last_edited_at = ?
		WHERE id = ?
	`,
		newContent,
		time.Now(),
		id.String(),
	)
	return err
}

func (r PostRepository) Archive(id domain.PostID) error {
	_, err := r.db.Exec(`
		UPDATE posts
		SET archived_at = ?
		WHERE id = ?
	`,
		time.Now(),
		id.String(),
	)
	return err
}

func dbPostToDomainPost(dbPost models.Post) *domain.Post {
	return domain.RebuildPost(
		domain.NewPostID(dbPost.ID),
		domain.NewUserID(dbPost.AuthorID),
		dbPost.Title,
		dbPost.Content,
		dbPost.CreatedAt,
		dbPost.LastEditedAt,
		dbPost.ArchivedAt,
	)
}

func dbPostsToDomainPosts(dbPosts []models.Post) []domain.Post {
	posts := []domain.Post{}
	for _, post := range dbPosts {
		posts = append(posts, *dbPostToDomainPost(post))
	}
	return posts
}
