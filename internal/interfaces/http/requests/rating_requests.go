package requests

import "blog/pkg/ddd/validation"

type CreateRatingRequest struct {
	PostID     string
	RatingType string
}

func (r CreateRatingRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.PostID, "post_id"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.Required(r.RatingType, "rating_type"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

type ChangeRatingRequest struct{}
