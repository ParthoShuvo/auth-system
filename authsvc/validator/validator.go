package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	validPasswdTag = "validPwd"
)

// PasswordValidator checks whether a password complies with minimal requirements.
type PasswordValidator struct {
	lc        *regexp.Regexp
	uc        *regexp.Regexp
	dc        *regexp.Regexp
	sc        *regexp.Regexp
	xc        *regexp.Regexp
	specChars string
}

// NewPasswordValidator creates a password validator.
func NewPasswordValidator() *PasswordValidator {
	const specChars = `_!@$%`
	return &PasswordValidator{
		lc:        regexp.MustCompile(`[a-z]`),
		uc:        regexp.MustCompile(`[A-Z]`),
		dc:        regexp.MustCompile(`[0-9]`),
		sc:        regexp.MustCompile(`[` + specChars + `]`),
		xc:        regexp.MustCompile(`[^a-zA-Z0-9` + specChars + `]`),
		specChars: specChars,
	}
}

// Validate a password.
func (pv *PasswordValidator) Validate(password string) error {
	if !pv.lc.MatchString(password) {
		return errors.New("1 lowercase letter")
	}
	if !pv.uc.MatchString(password) {
		return errors.New("1 uppercase letter")
	}
	if !pv.dc.MatchString(password) {
		return errors.New("1 digit")
	}
	if !pv.sc.MatchString(password) {
		return fmt.Errorf("1 special character (%s)", pv.specChars)
	}
	if illegalChars := pv.xc.FindAllString(password, -1); len(illegalChars) > 0 {
		return fmt.Errorf("password contains one or more illegal characters: %s", strings.Join(illegalChars, " "))
	}
	return nil
}

func New() *validator.Validate {
	validate := validator.New()
	registerCustomValidators(validate)
	return validate
}

func registerCustomValidators(validate *validator.Validate) {
	validate.RegisterValidation(validPasswdTag, passwordValidator)
}

func passwordValidator(fl validator.FieldLevel) bool {
	return NewPasswordValidator().Validate(fl.Field().String()) == nil
}
