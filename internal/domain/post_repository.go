package domain

type PostRepository interface {
	All() ([]Post, error)
	FindByID(id PostID) (*Post, error)
	Create(post *Post) (*Post, error)
	UpdateTitle(newTitle string) error
	UpdateContent(newContent string) error
	Archive(id PostID) error
}
