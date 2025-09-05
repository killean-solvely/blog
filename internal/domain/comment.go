package domain

import "blog/pkg/ddd"

type Comment struct {
	*ddd.AggregateBase
	commenterID UserID
	postID      PostID
	content     string
}
