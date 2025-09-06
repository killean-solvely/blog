package application

import (
	"time"

	"blog/internal/domain"
)

type PostDTO struct {
	AuthorID     string     `json:"author_id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	CreatedAt    time.Time  `json:"created_at"`
	LastEditedAt *time.Time `json:"last_edited_at"`
	ArchivedAt   *time.Time `json:"archived_at"`
}

func NewPostDTO(
	authorID, title, content string,
	createdAt time.Time,
	lastEditedAt, archivedAt *time.Time,
) *PostDTO {
	return &PostDTO{
		AuthorID:     authorID,
		Title:        title,
		Content:      content,
		CreatedAt:    createdAt,
		LastEditedAt: lastEditedAt,
		ArchivedAt:   archivedAt,
	}
}

func (dto *PostDTO) FromDomain(post *domain.Post) {
	dto = NewPostDTO(
		post.AuthorID().String(),
		post.Title(),
		post.Content(),
		post.CreatedAt(),
		post.LastEditedAt(),
		post.ArchivedAt(),
	)
}

func (dto PostDTO) ToDomain() *domain.Post {
	return domain.RebuildPost(
		domain.NewUserID(dto.AuthorID),
		dto.Title,
		dto.Content,
		dto.CreatedAt,
		dto.LastEditedAt,
		dto.ArchivedAt,
	)
}
