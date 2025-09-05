package domain

import "blog/pkg/ddd"

type Post struct {
	*ddd.AggregateBase
	authorID UserID
	title    string
	content  string
	likes    int
	dislikes int
}
