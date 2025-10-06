package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/validate"
	"go.uber.org/zap"
)

// Validator wraps gookit/validate for structured validation
type Validator struct {
	validate *validator.Validate
	logger   *log.Logger
}

// NewValidator creates a new validator instance
func NewValidator(logger *log.Logger) *Validator {
	return &Validator{
		logger: logger,
	}
}

// ValidateStruct validates a struct using validation tags
func (v *Validator) ValidateStruct(data interface{}) error {
	validator := validate.Struct(data)
	
	if !validator.Validate() {
		errors := validator.Errors
		var errorMessages []string
		
		for field, messages := range errors {
			for _, message := range messages {
				errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", field, message))
			}
		}
		
		v.logger.Debug("Validation failed", 
			zap.Strings("errors", errorMessages),
		)
		
		return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}
	
	return nil
}

// ValidateEmail validates an email address
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// ValidatePassword validates a password
func (v *Validator) ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	
	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Check for at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	
	if !hasUpper || !hasLower || !hasDigit {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, and one digit")
	}
	
	return nil
}

// ValidateUsername validates a username
func (v *Validator) ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return fmt.Errorf("username must be between 3 and 50 characters")
	}
	
	// Username should only contain alphanumeric characters and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	
	return nil
}