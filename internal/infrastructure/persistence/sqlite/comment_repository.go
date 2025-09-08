package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"blog/internal/domain"
	"blog/internal/infrastructure/persistence/models"

	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r CommentRepository) All() ([]domain.Comment, error) {
	var dbComments []models.Comment
	err := r.db.Select(&dbComments, "SELECT * FROM comments")
	if err != nil {
		return nil, err
	}

	comments := dbCommentsToDomainComments(dbComments)
	return comments, nil
}

func (r CommentRepository) FindByID(id domain.CommentID) (*domain.Comment, error) {
	var dbComment models.Comment
	err := r.db.Get(&dbComment, "SELECT * FROM comments WHERE id=?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCommentNotFound
		}
		return nil, err
	}

	comment := dbCommentToDomainComment(dbComment)
	return comment, nil
}

func (r *CommentRepository) FindByUser(userID domain.UserID) ([]domain.Comment, error) {
	var dbComments []models.Comment
	err := r.db.Select(&dbComments, "SELECT * FROM comments WHERE commenter_id=?", userID)
	if err != nil {
		return nil, err
	}

	comments := dbCommentsToDomainComments(dbComments)
	return comments, nil
}

func (r *CommentRepository) FindByPost(postID domain.PostID) ([]domain.Comment, error) {
	var dbComments []models.Comment
	err := r.db.Select(&dbComments, "SELECT * FROM comments WHERE post_id=?", postID)
	if err != nil {
		return nil, err
	}

	comments := dbCommentsToDomainComments(dbComments)
	return comments, nil
}

func (r *CommentRepository) Exists(id domain.CommentID) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM comments WHERE id=?", id)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *CommentRepository) Create(comment *domain.Comment) (*domain.Comment, error) {
	_, err := r.db.Exec(`
		INSERT INTO 
		comments (id, post_id, commenter_id, content, created_at, last_updated_at, archived_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		comment.GetID().String(),
		comment.PostID().String(),
		comment.CommenterID().String(),
		comment.Content(),
		comment.CreatedAt(),
		comment.LastUpdatedAt(),
		comment.ArchivedAt(),
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepository) UpdateContent(id domain.CommentID, newContent string) error {
	_, err := r.db.Exec(`
		UPDATE comments
		SET content = ?, last_updated_at = ?
		WHERE id = ?
	`,
		newContent,
		time.Now(),
		id.String(),
	)
	return err
}

func (r *CommentRepository) Archive(id domain.CommentID) error {
	_, err := r.db.Exec(`
		UPDATE comments
		SET archived_at = ?
		WHERE id = ?
	`,
		time.Now(),
		id.String(),
	)
	return err
}

func dbCommentToDomainComment(dbComment models.Comment) *domain.Comment {
	return domain.RebuildComment(
		domain.NewCommentID(dbComment.ID),
		domain.NewPostID(dbComment.PostID),
		domain.NewUserID(dbComment.CommenterID),
		dbComment.Content,
		dbComment.CreatedAt,
		dbComment.LastUpdatedAt,
		dbComment.ArchivedAt,
	)
}

func dbCommentsToDomainComments(dbComments []models.Comment) []domain.Comment {
	comments := []domain.Comment{}
	for _, comment := range dbComments {
		comments = append(comments, *dbCommentToDomainComment(comment))
	}
	return comments
}
