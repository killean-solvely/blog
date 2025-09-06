package domain

import "errors"

var (
	// User
	ErrDescriptionTooLong = errors.New("description cannot exceed 255 character limit")
	ErrMissingUserRoles   = errors.New("cannot create user without a role")

	// Post
	ErrTitleCannotBeEmpty   = errors.New("post title cannot be empty")
	ErrContentCannotBeEmpty = errors.New("post content cannot be empty")

	// Comment
	ErrCommentCannotBeEmpty = errors.New("comment cannot be empty")
)
