package cnab

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"time"

	ourErrors "github.com/RafOSS-br/cnab-stream/errors"
	iError "github.com/RafOSS-br/cnab-stream/internal/error"
)

// Integer Field Handlers

// validateIntField validates an integer field.
func validateIntField(field FieldSpec) error {
	// No specific validation needed for integer fields
	return nil
}

// parseIntField parses an integer field.
func parseIntField(field FieldSpec, rawValue []byte) (interface{}, error) {
	trimmedValue := bytes.TrimSpace(rawValue)
	if len(trimmedValue) == 0 {
		return nil, ourErrors.CNAB_ErrFieldIsEmpty.Creator(fieldToError("Name", field.Name))
	}
	return atoiUnsafe(trimmedValue)
}

// formatIntField formats an integer field.
func formatIntField(field FieldSpec, value interface{}) (string, error) {
	intValue, err := toInt(value)
	if err != nil {
		return "", ourErrors.CNAB_ErrFieldValueIsNotAnInt.Creator(fmt.Errorf("field %s: %w", field.Name, err))
	}
	return fmt.Sprintf("%0*d", field.Length, intValue), nil
}

// Float Field Handlers

var (
	// ErrInvalidDecimalValue is returned when the decimal value is negative
	ErrInvalidDecimalValue   = ourErrors.CNAB_ErrInvalidDecimalValue.Creator(fieldToError("Name", "decimal"))
	IsErrInvalidDecimalValue = iError.MatchError(ErrInvalidDecimalValue)
)

// validateFloatField validates a float field.
func validateFloatField(field FieldSpec) error {
	if field.Decimal < 0 {
		return ourErrors.CNAB_ErrInvalidDecimalValue.Creator(fieldToError("Name", field.Name))
	}
	return nil
}

// parseFloatField parses a float field.
func parseFloatField(field FieldSpec, rawValue []byte) (interface{}, error) {
	trimmedValue := bytes.TrimSpace(rawValue)
	if len(trimmedValue) == 0 {
		return nil, ourErrors.CNAB_ErrFieldIsEmpty.Creator(fieldToError("Name", field.Name))
	}
	return parseFloatBytes(trimmedValue, field.Decimal)
}

// formatFloatField formats a float field.
func formatFloatField(field FieldSpec, value interface{}) (string, error) {
	floatValue, err := toFloat(value)
	if err != nil {
		return "", ourErrors.CNAB_ErrFieldValueIsNotAnFloat.Creator(fieldToError("Name", field.Name))
	}
	scaledValue := int64(floatValue * pow10(field.Decimal))
	return fmt.Sprintf("%0*d", field.Length, scaledValue), nil
}

// Date Field Handlers

var (
	// ErrMissingDateFormat is returned when the date format is missing
	ErrMissingDateFormat   = ourErrors.CNAB_ErrMissingDateFormat.Creator(fieldToError("Name", "format"))
	IsErrMissingDateFormat = iError.MatchError(ErrMissingDateFormat)
)

// validateDateField validates a date field.
func validateDateField(field FieldSpec) error {
	if field.Format == "" {
		return ourErrors.CNAB_ErrMissingDateFormat.Creator(fieldToError("Name", field.Name))
	}
	return nil
}

// parseDateField parses a date field.
func parseDateField(field FieldSpec, rawValue []byte) (interface{}, error) {
	trimmedValue := bytes.TrimSpace(rawValue)
	if len(trimmedValue) == 0 {
		return nil, ourErrors.CNAB_ErrFieldIsEmpty.Creator(fieldToError("Name", field.Name))
	}
	if len(trimmedValue) != len(field.Format) {
		return nil, ourErrors.CNAB_ErrInvalidDateLength.Creator(fieldToError("Name", field.Name))
	}
	return parseDateBytes(trimmedValue, field.Format)
}

// formatDateField formats a date field.
func formatDateField(field FieldSpec, value interface{}) (string, error) {
	dateValue, ok := value.(time.Time)
	if !ok {
		return "", ourErrors.CNAB_ErrFieldValueIsNotAnDate.Creator(fieldToError("Name", field.Name))
	}
	return formatDate(dateValue, field.Format), nil
}

// String Field Handlers

// validateStringField validates a string field.
func validateStringField(field FieldSpec) error {
	// No specific validation needed for string fields
	return nil
}

// parseStringField parses a string field.
func parseStringField(field FieldSpec, rawValue []byte) (interface{}, error) {
	return string(bytes.TrimSpace(rawValue)), nil
}

// formatStringField formats a string field.
func formatStringField(field FieldSpec, value interface{}) (string, error) {
	strValue, ok := value.(string)
	if !ok {
		return "", ourErrors.CNAB_ErrFieldValueIsNotAnString.Creator(fieldToError("Name", field.Name))
	}
	return strValue, nil
}

// Utility Functions

// atoiUnsafe is a faster version of strconv.Atoi that does not check for errors.
func atoiUnsafe(b []byte) (int, error) {
	n := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, ourErrors.CNAB_ErrFieldValueIsNotAnInt.Creator(fmt.Errorf("invalid character %c", c))
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// parseFloatBytes parses a float from a byte slice.
func parseFloatBytes(b []byte, decimal int) (float64, error) {
	intValue, err := atoiUnsafe(b)
	if err != nil {
		return 0, ourErrors.CNAB_ErrFieldValueIsNotAnFloat.Creator(err)
	}
	return float64(intValue) / pow10(decimal), nil
}

var (
	// ErrInvalidDateLength is returned when the date length is invalid
	ErrInvalidDateLength   = ourErrors.CNAB_ErrInvalidDateLength.Err
	IsErrInvalidDateLength = iError.MatchError(ErrInvalidDateLength)
)

// parseDateBytes parses a date from a byte slice using a CNAB date format.
func parseDateBytes(b []byte, format string) (time.Time, error) {
	if len(b) != len(format) {
		return time.Time{}, ErrInvalidDateLength
	}
	goFormat := convertDateFormat(format)
	timeValue, err := time.Parse(goFormat, string(b))
	if err != nil {
		return time.Time{}, ourErrors.CNAB_ErrFieldValueIsNotAnDate.Creator(err)
	}
	return timeValue, nil
}

// pow10 returns 10^n
func pow10(n int) float64 {
	result := 1.0
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

// convertDateFormat converts a CNAB date format to a Go date format.
func convertDateFormat(format string) string {
	replacements := map[string]string{
		"YYYY": "2006",
		"MM":   "01",
		"DD":   "02",
	}
	for old, new := range replacements {
		format = strings.ReplaceAll(format, old, new)
	}
	return format
}

// formatDate formats a time.Time value using a CNAB date format.
func formatDate(value time.Time, format string) string {
	goFormat := convertDateFormat(format)
	return value.Format(goFormat)
}

// toInt converts a value to an int.
func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return 0, ourErrors.CNAB_ErrFieldValueIsNotAnInt.Creator(fieldToError("Value", v))
		}
		return intValue, nil
	default:
		return 0, ourErrors.CNAB_ErrFieldValueIsNotAnInt.Creator(fieldToError("Value", fmt.Sprintf("%v", value)))
	}
}

// toFloat converts a value to a float64.
func toFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		floatValue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, ourErrors.CNAB_ErrFieldValueIsNotAnFloat.Creator(fieldToError("Value", v))
		}
		return floatValue, nil
	default:
		return 0, ourErrors.CNAB_ErrFieldValueIsNotAnFloat.Creator(fieldToError("Value", fmt.Sprintf("%v", value)))
	}
}

// fieldToError returns an error message for a field.
func fieldToError(fieldName, fieldValue string) error {
	return fmt.Errorf("field %s %s", fieldName, fieldValue)
}
