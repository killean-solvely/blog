package domain

type CommentID string

func NewCommentID(id string) CommentID {
	return CommentID(id)
}

func (id CommentID) String() string {
	return string(id)
}
