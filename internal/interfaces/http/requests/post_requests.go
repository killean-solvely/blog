package requests

import "blog/pkg/ddd/validation"

type CreatePostRequest struct {
	AuthorID string `json:"author_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

func (r CreatePostRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.AuthorID, "author_id"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.Required(r.Title, "title"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.Required(r.Content, "content"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

type UpdatePostTitle struct {
	Title string `json:"title"`
}

func (r UpdatePostTitle) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.Title, "title"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

type UpdatePostContent struct {
	Content string `json:"content"`
}

func (r UpdatePostContent) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.Content, "content"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
