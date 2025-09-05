package domain

type PostID string

func NewPostID(id string) PostID {
	return PostID(id)
}

func (id PostID) String() string {
	return string(id)
}
