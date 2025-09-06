package domain

type RatingType string

const (
	RatingTypeLike    RatingType = "like"
	RatingTypeDislike RatingType = "dislike"
)

func (rt RatingType) String() string {
	return string(rt)
}
