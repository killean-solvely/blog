package sqlite

import (
	"database/sql"
	"errors"

	"blog/internal/domain"
	"blog/internal/infrastructure/persistence/models"

	"github.com/jmoiron/sqlx"
)

type RatingRepository struct {
	db *sqlx.DB
}

func NewRatingRepository(db *sqlx.DB) *RatingRepository {
	return &RatingRepository{
		db: db,
	}
}

func (r RatingRepository) All() ([]domain.Rating, error) {
	var dbRatings []models.Rating
	err := r.db.Select(&dbRatings, "SELECT * FROM ratings")
	if err != nil {
		return nil, err
	}

	ratings := dbRatingsToDomainRatings(dbRatings)
	return ratings, nil
}

func (r RatingRepository) FindByID(id domain.RatingID) (*domain.Rating, error) {
	var dbRating models.Rating
	err := r.db.Get(&dbRating, "SELECT * FROM ratings WHERE id=?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRatingNotFound
		}
		return nil, err
	}

	rating := dbRatingToDomainRating(dbRating)
	return rating, nil
}

func (r *RatingRepository) FindByUser(userID domain.UserID) ([]domain.Rating, error) {
	var dbRatings []models.Rating
	err := r.db.Select(&dbRatings, "SELECT * FROM ratings WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	ratings := dbRatingsToDomainRatings(dbRatings)
	return ratings, nil
}

func (r *RatingRepository) FindByPost(postID domain.PostID) ([]domain.Rating, error) {
	var dbRatings []models.Rating
	err := r.db.Select(&dbRatings, "SELECT * FROM ratings WHERE post_id=?", postID)
	if err != nil {
		return nil, err
	}

	ratings := dbRatingsToDomainRatings(dbRatings)
	return ratings, nil
}

func (r *RatingRepository) Exists(id domain.RatingID) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM ratings WHERE id=?", id)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RatingRepository) ExistsOnPostByUser(
	postID domain.PostID,
	userID domain.UserID,
) (bool, error) {
	var count int
	err := r.db.Get(
		&count,
		"SELECT COUNT(*) FROM ratings WHERE post_id=? AND user_id=?",
		postID,
		userID,
	)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RatingRepository) Create(rating *domain.Rating) (*domain.Rating, error) {
	_, err := r.db.Exec(`
		INSERT INTO 
		ratings (id, post_id, user_id, rating_type, created_at) 
		VALUES (?, ?, ?, ?, ?)
	`,
		rating.GetID().String(),
		rating.PostID().String(),
		rating.UserID().String(),
		string(rating.RatingType()),
		rating.CreatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return rating, nil
}

func (r *RatingRepository) ChangeRating(id domain.RatingID, newRatingType domain.RatingType) error {
	_, err := r.db.Exec(`
		UPDATE ratings
		SET rating_type = ?
		WHERE id = ?
	`,
		string(newRatingType),
		id.String(),
	)
	return err
}

func (r *RatingRepository) RemoveRating(id domain.RatingID) error {
	_, err := r.db.Exec(`
		DELETE FROM ratings
		WHERE id = ?
	`,
		id.String(),
	)
	return err
}

func dbRatingToDomainRating(dbRating models.Rating) *domain.Rating {
	return domain.RebuildRating(
		domain.NewRatingID(dbRating.ID),
		domain.NewPostID(dbRating.PostID),
		domain.NewUserID(dbRating.UserID),
		domain.RatingType(dbRating.RatingType),
		dbRating.CreatedAt,
		dbRating.UpdatedAt,
	)
}

func dbRatingsToDomainRatings(dbRatings []models.Rating) []domain.Rating {
	ratings := []domain.Rating{}
	for _, rating := range dbRatings {
		ratings = append(ratings, *dbRatingToDomainRating(rating))
	}
	return ratings
}

