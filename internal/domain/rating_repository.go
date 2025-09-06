package domain

type RatingRepository interface {
	All() ([]Rating, error)
	FindByID(id RatingID) (*Rating, error)
	FindByUser(userID UserID) ([]Rating, error)
	FindByPost(postID PostID) ([]Rating, error)
	Create(rating *Rating) (*Rating, error)
	ChangeRating(id RatingID, newRatingType RatingType) error
	RemoveRating(id RatingID) error
}
