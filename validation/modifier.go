package validation

import (
	"context"
	"reflect"

	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
)

type MoldModValidator struct {
	validator *validator.Validate
}

func NewMoldModValidator() *MoldModValidator {
	v := validator.New()
	// Set the tag name to "binding", SetTagName allows for changing of the default tag name of 'validate'
	v.SetTagName("binding")
	// Register Tag Name Function to get json name as alternate names for StructFields.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("json")
	})
	// Register custom validation tags if needed
	v.RegisterValidation("customValidation", customValidationFunc)
	return &MoldModValidator{validator: v}
}

// ValidateStruct is called by Gin to validate the struct
func (cv *MoldModValidator) ValidateStruct(obj interface{}) error {
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
func (cv *MoldModValidator) Engine() interface{} {
	return cv.validator
}

// Custom validation function
func customValidationFunc(fl validator.FieldLevel) bool {
	// Custom validation logic here
	return true
}
