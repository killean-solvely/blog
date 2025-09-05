package domain

import "errors"

var ErrDescriptionTooLong = errors.New("description cannot exceed 255 character limit")
