package domain

import "errors"

var (
	// Comment
	ErrCommentNotFound      = errors.New("comment not found")
	ErrCommentCannotBeEmpty = errors.New("comment cannot be empty")

	// Post
	ErrPostNotFound         = errors.New("post not found")
	ErrTitleCannotBeEmpty   = errors.New("post title cannot be empty")
	ErrContentCannotBeEmpty = errors.New("post content cannot be empty")

	// Rating
	ErrRatingNotFound = errors.New("rating now found")

	// User
	ErrUserNotFound       = errors.New("user not found")
	ErrDescriptionTooLong = errors.New("description cannot exceed 255 character limit")
	ErrMissingUserRoles   = errors.New("cannot create user without a role")
)
