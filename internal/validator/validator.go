package validator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	// use a single instance of Validate, it caches struct info
	Validate      *validator.Validate
	EchoValidator *CustomEchoValidator
)

func init() {
	Validate = validator.New()
	EchoValidator = &CustomEchoValidator{validator: Validate}
}

type CustomEchoValidator struct {
	validator *validator.Validate
}

func (cev *CustomEchoValidator) Validate(i any) error {
	if err := cev.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
