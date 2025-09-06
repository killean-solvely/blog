package application

import (
	"time"

	"blog/internal/domain"
)

type PostDTO struct {
	ID           string     `json:"id"`
	AuthorID     string     `json:"author_id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	CreatedAt    time.Time  `json:"created_at"`
	LastEditedAt *time.Time `json:"last_edited_at"`
	ArchivedAt   *time.Time `json:"archived_at"`
}

func NewPostDTO(
	id, authorID, title, content string,
	createdAt time.Time,
	lastEditedAt, archivedAt *time.Time,
) *PostDTO {
	return &PostDTO{
		ID:           id,
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
		post.GetID().String(),
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
		domain.NewPostID(dto.ID),
		domain.NewUserID(dto.AuthorID),
		dto.Title,
		dto.Content,
		dto.CreatedAt,
		dto.LastEditedAt,
		dto.ArchivedAt,
	)
}

type CommentDTO struct {
	ID            string
	PostID        string
	CommenterID   string
	Content       string
	CreatedAt     time.Time
	LastUpdatedAt *time.Time
	ArchivedAt    *time.Time
}

func NewCommentDTO(
	id, postID, commenterID, content string,
	createdAt time.Time,
	lastUpdatedAt *time.Time,
	archivedAt *time.Time,
) *CommentDTO {
	return &CommentDTO{
		ID:            id,
		PostID:        postID,
		CommenterID:   commenterID,
		Content:       content,
		CreatedAt:     createdAt,
		LastUpdatedAt: lastUpdatedAt,
		ArchivedAt:    archivedAt,
	}
}

func (dto *CommentDTO) FromDomain(comment *domain.Comment) {
	dto = NewCommentDTO(
		comment.GetID().String(),
		comment.PostID().String(),
		comment.CommenterID().String(),
		comment.Content(),
		comment.CreatedAt(),
		comment.LastUpdatedAt(),
		comment.ArchivedAt(),
	)
}

func (dto CommentDTO) ToDomain() *domain.Comment {
	return domain.RebuildComment(
		domain.NewCommentID(dto.ID),
		domain.NewPostID(dto.PostID),
		domain.NewUserID(dto.CommenterID),
		dto.Content,
		dto.CreatedAt,
		dto.LastUpdatedAt,
		dto.ArchivedAt,
	)
}

type UserDTO struct {
	ID           string
	Email        string
	PasswordHash string
	Username     string
	Description  string
	UserRoles    []string
	JoinDate     time.Time
}

func NewUserDTO(
	id, email, username, passwordHash, description string,
	userRoles []string,
	joinDate time.Time,
) *UserDTO {
	return &UserDTO{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
		Description:  description,
		UserRoles:    userRoles,
		JoinDate:     joinDate,
	}
}

func (dto *UserDTO) FromDomain(user *domain.User) {
	roles := []string{}
	for _, role := range user.UserRoles() {
		roles = append(roles, role.String())
	}

	dto = NewUserDTO(
		user.GetID().String(),
		user.Email(),
		user.Username(),
		user.PasswordHash(),
		user.Description(),
		roles,
		user.JoinDate(),
	)
}

func (dto *UserDTO) ToDomain() *domain.User {
	roles := []domain.UserRole{}
	for _, role := range dto.UserRoles {
		roles = append(roles, domain.UserRole(role))
	}

	return domain.RebuildUser(
		domain.NewUserID(dto.ID),
		dto.Email,
		dto.PasswordHash,
		dto.Username,
		dto.Description,
		roles,
		dto.JoinDate,
	)
}
