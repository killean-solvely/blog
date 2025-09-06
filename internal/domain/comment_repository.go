package domain

type CommentRepository interface {
	All() ([]Comment, error)
	FindByID(id CommentID) (*Comment, error)
	FindByUser(userID UserID) ([]Comment, error)
	FindByPost(postID PostID) ([]Comment, error)
	Create(comment *Comment) (*Comment, error)
	UpdateContent(newContent string) error
	Archive(id CommentID) error
}
