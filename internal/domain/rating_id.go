package domain

type RatingID string

func NewRatingID(id string) RatingID {
	return RatingID(id)
}

func (id RatingID) String() string {
	return string(id)
}
