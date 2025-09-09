package requests

import "blog/pkg/ddd/validation"

type CreateCommentRequest struct {
	Content string `json:"content"`
}

func (r CreateCommentRequest) Validate() *validation.Errors {
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

type EditCommentRequest struct {
	Content string `json:"content"`
}

func (r EditCommentRequest) Validate() *validation.Errors {
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