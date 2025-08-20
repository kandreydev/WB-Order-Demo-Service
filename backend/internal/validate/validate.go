package validate

import (
	"github.com/GkadyrG/L0/backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var validate = validator.New()

func ValidateOrder(order model.Order) error {
	_ = validate.RegisterValidation("e164", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if len(phone) < 4 || phone[0] != '+' {
			return false
		}
		return true
	})

	if err := validate.Struct(order); err != nil {
		for _, verr := range err.(validator.ValidationErrors) {
			return errors.Errorf("validation failed for field '%s', tag '%s'", verr.StructNamespace(), verr.Tag())
		}
	}

	return nil
}
