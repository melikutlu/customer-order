package internal

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"strings"
	"tesodev-korpes/CustomerService/internal/types"
)

func ToCustomerResponse(customer *types.Customer) *types.CustomerResponseModel {
	return &types.CustomerResponseModel{
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Username:  customer.Username,
		Email:     customer.Email,
	}
}

func ToCustomerRespList(customers []types.Customer) []types.CustomerResponseModel {
	var customerResponseModels []types.CustomerResponseModel
	for _, customer := range customers {
		customerResponseModels = append(customerResponseModels, types.CustomerResponseModel{
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Username:  customer.Username,
			Email:     customer.Email,
		})
	}
	return customerResponseModels
}

func ToLoginResponse(customer *types.Customer) *types.CustomerLoginResponseModel {
	return &types.CustomerLoginResponseModel{
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Username:  customer.Username,
		Email:     customer.Email,
		Id:        customer.Id,
		Password:  customer.Password,
	}
}

func ValidateEmail(r *types.CustomerRequestModel) error {

	email := r.Email

	if email == "" {
		return errors.New("Email is required")
	}

	if !strings.Contains(email, "@") {
		return errors.New("Email must contain @")
	}
	return nil

}

func ValidateAge(r *types.CustomerRequestModel) error {
	age := r.Age
	if age == 0 {
		return errors.New("Age is required")
	}

	if age < 18 {
		return errors.New("Age must be 18 or older")
	}

	return nil
}

func ValidateFirstLetterUpperCase(customer *types.CustomerRequestModel) error {
	errors := make(map[string]string)
	if customer.FirstName != "" {
		if !isFirstLetterUpperCase(customer.FirstName) {
			errors["FirstName"] = "First name must start with an uppercase letter"
		}
		if containsDigit(customer.FirstName) {
			errors["FirstName"] = "First name contains a number"
		}
	}
	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}
	return nil
}

func containsDigit(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return true
		}
	}
	return false
}

func isFirstLetterUpperCase(s string) bool {
	if len(s) > 0 {
		return strings.ToUpper(s[:1]) == s[:1]
	}
	return false
}

func ValidateCustomer(customer *types.CustomerRequestModel, validate *validator.Validate) error {
	validationErrors := make(map[string]string)

	if err := ValidateAge(customer); err != nil {
		validationErrors["Age"] = err.Error()
	}

	if err := ValidateEmail(customer); err != nil {
		validationErrors["Email"] = err.Error()
	}

	if err := ValidateFirstLetterUpperCase(customer); err != nil {
		// Use the errors from ValidateFirstLetterUpperCase directly
		if valErr, ok := err.(*ValidationError); ok {
			for field, msg := range valErr.Errors {
				validationErrors[field] = msg
			}
		}
	}

	if err := validate.Struct(customer); err != nil {
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range fieldErrors {
				if fieldError.Tag() == "required" {
					validationErrors[fieldError.Field()] = "This field is required"
				}

			}
		}
	}

	if len(validationErrors) > 0 {
		return &ValidationError{Errors: validationErrors}
	}

	return nil
}

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "Validation failed "
}
