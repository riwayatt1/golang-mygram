package utils

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}

func IsGinH(data interface{}) bool {
	// Get the reflect.Type of the data
	dataType := reflect.TypeOf(data)

	// Check if the data type is gin.H
	ginHType := reflect.TypeOf(gin.H{})
	if dataType == ginHType {
		return true
	}

	// If not gin.H, check if it's a map[string]interface{}
	if dataType.Kind() == reflect.Map {
		// Check if the map has string keys and interface{} values
		keyType := dataType.Key()
		valueType := dataType.Elem()
		if keyType.Kind() == reflect.String && valueType.Kind() == reflect.Interface {
			return true
		}
	}

	return false
}

// RemoveFieldsFromMap removes specified fields from a struct and returns a map with lowercase keys and underscores
func RemoveFieldsFromMap(data interface{}, fields ...string) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		fieldName := field.Name

		// Skip fields specified in the fields parameter
		if contains(fields, fieldName) {
			continue
		}

		// Convert fieldName from camelCase to lowercase with underscores, except for "id"
		key := camelToSnake(fieldName)

		// Special case for "id" to keep it lowercase
		if fieldName == "ID" {
			key = "id"
		}

		// Get field value
		fieldValue := val.Field(i).Interface()
		result[key] = fieldValue
	}

	return result
}

// contains checks if slice contains string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// camelToSnake converts camelCase string to snake_case
func camelToSnake(s string) string {
	var output []rune
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			output = append(output, '_')
		}
		output = append(output, r)
	}
	return strings.ToLower(string(output))
}
