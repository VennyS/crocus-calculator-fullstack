package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
	Validate.RegisterValidation("not_pending", func(fl validator.FieldLevel) bool {
		return fl.Field().String() != "pending"
	})
}

func ValidateStruct(s any) error {
	if Validate == nil {
		InitValidator()
	}

	err := Validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var multiErr error
	for _, vErr := range validationErrors {
		var errMsg string
		switch vErr.Tag() {
		case "required":
			errMsg = fmt.Sprintf("field '%s' is required", strings.ToLower(vErr.Field()))
		case "oneof":
			errMsg = fmt.Sprintf("field '%s' must be one of: %s",
				strings.ToLower(vErr.Field()),
				strings.Join(strings.Split(vErr.Param(), " "), ", "))
		default:
			errMsg = fmt.Sprintf("field '%s' failed '%s' validation",
				strings.ToLower(vErr.Field()),
				vErr.Tag())
		}

		if multiErr == nil {
			multiErr = errors.New(errMsg)
		} else {
			multiErr = fmt.Errorf("%w; %s", multiErr, errMsg)
		}
	}

	return multiErr
}
