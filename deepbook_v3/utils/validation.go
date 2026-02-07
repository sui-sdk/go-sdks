package utils

import (
	"fmt"

	suiutils "github.com/sui-sdks/go-sdks/sui/utils"
)

func ValidateRequired[T any](value *T, errorMessage string) T {
	if value == nil {
		panic(&ConfigurationError{DeepBookError{Msg: errorMessage}})
	}
	return *value
}

func ValidateAddress(address string, fieldName string) string {
	if fieldName == "" {
		fieldName = "Address"
	}
	if !suiutils.IsValidSuiAddress(address) {
		panic(&ValidationError{DeepBookError{Msg: fmt.Sprintf("%s must be a valid Sui address", fieldName)}})
	}
	return address
}

func ValidatePositiveNumber(value float64, fieldName string) float64 {
	if value <= 0 {
		panic(&ValidationError{DeepBookError{Msg: fmt.Sprintf("%s must be a positive number", fieldName)}})
	}
	return value
}

func ValidateNonNegativeNumber(value float64, fieldName string) float64 {
	if value < 0 {
		panic(&ValidationError{DeepBookError{Msg: fmt.Sprintf("%s must be non-negative", fieldName)}})
	}
	return value
}

func ValidateRange(value, min, max float64, fieldName string) float64 {
	if value < min || value > max {
		panic(&ValidationError{DeepBookError{Msg: fmt.Sprintf("%s must be between %v and %v", fieldName, min, max)}})
	}
	return value
}

func ValidateNonEmptyArray[T any](array []T, fieldName string) []T {
	if len(array) == 0 {
		panic(&ValidationError{DeepBookError{Msg: fmt.Sprintf("%s cannot be empty", fieldName)}})
	}
	return array
}
