package validators

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type CustomValidator struct {
	validator *validator.Validate
}

type ValidationError struct {
	Namespace string `json:"namespace,omitempty"`
	Field     string `json:"field,omitempty"`
	Error     string `json:"error,omitempty"`
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var errs []ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, ValidationError{
				Namespace: err.Namespace(),
				Field:     err.Field(),
				Error:     fmt.Sprintf("%s - %s", err.Type(), err.Tag()),
			})
		}
		return echo.NewHTTPError(http.StatusBadRequest, errs)
	}
	return nil

}

func NewValidator() *CustomValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	v.RegisterValidation("interval", intervalValidator)
	return &CustomValidator{validator: v}
}

func intervalValidator(fl validator.FieldLevel) bool {
	interval := fl.Field().String()
	parts := strings.Split(interval, "-")
	if len(parts) != 2 {
		return false
	}
	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}
	if start < 9 || end > 21 {
		return false
	}
	return true
}
