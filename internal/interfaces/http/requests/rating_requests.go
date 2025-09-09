package requests

import "blog/pkg/ddd/validation"

type CreateRatingRequest struct {
	PostID     string `json:"post_id"`
	RatingType string `json:"rating_type"`
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

type ChangeRatingRequest struct {
	RatingID   string `json:"rating_id"`
	PostID     string `json:"post_id"`
	RatingType string `json:"rating_type"`
}

func (r ChangeRatingRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.RatingID, "rating_id"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

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

type RemoveRatingRequest struct {
	RatingID string `json:"rating_id"`
}

func (r RemoveRatingRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.RatingID, "rating_id"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
