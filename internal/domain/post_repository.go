package domain

type PostRepository interface {
	All() ([]Post, error)
	FindByID(id PostID) (*Post, error)
	FindByAuthor(authorID UserID) ([]Post, error)
	Exists(id PostID) (bool, error)
	Create(post *Post) (*Post, error)
	UpdateTitle(id PostID, newTitle string) error
	UpdateContent(id PostID, newContent string) error
	Archive(id PostID) error
}
