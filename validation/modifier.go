package validation

import (
	"context"

	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
)

type MoldModifyValidator struct {
	validator *validator.Validate
}

func NewMoldModifyValidator(v *validator.Validate) *MoldModifyValidator {
	return &MoldModifyValidator{validator: v}
}

// ValidateStruct is called by Gin to validate the struct
func (cv *MoldModifyValidator) ValidateStruct(obj interface{}) error {
	// transform the object using mold before validating the struct
	transformer := modifiers.New()
	if err := transformer.Struct(context.Background(), obj); err != nil {
		return err
	}
	// validate the struct
	if err := cv.validator.Struct(obj); err != nil {
		return err
	}
	return nil
}

// Engine is called by Gin to retrieve the underlying validation engine
func (cv *MoldModifyValidator) Engine() interface{} {
	return cv.validator
}
