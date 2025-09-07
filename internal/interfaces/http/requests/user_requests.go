package requests

import "blog/pkg/ddd/validation"

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RegisterUserRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.Email, "email"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.Email(r.Email, "email"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.Required(r.Username, "username"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.MinLength(r.Username, "username", 3); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.Required(r.Password, "password"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.MinLength(r.Password, "password", 8); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r LoginUserRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.Username, "username"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.Required(r.Password, "password"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

type SetUserRolesRequest struct {
	UserRoles []string `json:"user_roles"`
}

func (r SetUserRolesRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.RequiredStringSlice(r.UserRoles, "user_roles"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.MinLengthStringSlice(r.UserRoles, "user_roles", 1); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

type UpdateUserDescriptionRequest struct {
	Description string `json:"description"`
}

func (r UpdateUserDescriptionRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.Description, "description"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password"`
}

func (r UpdateUserPasswordRequest) Validate() *validation.Errors {
	v := validation.New()
	errors := validation.NewErrors()

	if err := v.Required(r.Password, "password"); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if err := v.MinLength(r.Password, "password", 8); err != nil {
		errors.ValidationErrors = append(errors.ValidationErrors, *err)
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}
