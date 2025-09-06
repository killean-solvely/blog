package domain

type RatingRepository interface {
	All() ([]Rating, error)
	FindByID(id RatingID) (*Rating, error)
	FindByUser(userID UserID) ([]Rating, error)
	FindByPost(postID PostID) ([]Rating, error)
	Exists(id RatingID) (bool, error)
	ExistsOnPostByUser(postID PostID, userID UserID) (bool, error)
	Create(rating *Rating) (*Rating, error)
	ChangeRating(id RatingID, newRatingType RatingType) error
	RemoveRating(id RatingID) error
}
