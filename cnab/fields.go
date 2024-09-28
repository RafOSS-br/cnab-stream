package cnab

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"time"

	ourErrors "github.com/RafOSS-br/cnab-stream/errors"
	iError "github.com/RafOSS-br/cnab-stream/internal/error"
)

// Integer Field Handlers

func validateIntField(field FieldSpec) error {
	// No specific validation needed for integer fields
	return nil
}

func parseIntField(field FieldSpec, rawValue []byte) (interface{}, error) {
	trimmedValue := bytes.TrimSpace(rawValue)
	if len(trimmedValue) == 0 {
		return nil, ourErrors.CNAB_ErrFieldIsEmpty.Creator(fieldToError("Name", field.Name))
	}
	return atoiUnsafe(trimmedValue)
}

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

func validateFloatField(field FieldSpec) error {
	if field.Decimal < 0 {
		return ourErrors.CNAB_ErrInvalidDecimalValue.Creator(fieldToError("Name", field.Name))
	}
	return nil
}

func parseFloatField(field FieldSpec, rawValue []byte) (interface{}, error) {
	trimmedValue := bytes.TrimSpace(rawValue)
	if len(trimmedValue) == 0 {
		return nil, ourErrors.CNAB_ErrFieldIsEmpty.Creator(fieldToError("Name", field.Name))
	}
	return parseFloatBytes(trimmedValue, field.Decimal)
}

func formatFloatField(field FieldSpec, value interface{}) (string, error) {
	floatValue, err := toFloat(value)
	if err != nil {
		return "", ourErrors.CNAB_ErrFieldValueIsNotAnFloat.Creator(fieldToError("Name", field.Name))
	}
	scaledValue := int64(floatValue * pow10(field.Decimal))
	return fmt.Sprintf("%0*d", field.Length, scaledValue), nil
}

// Date Field Handlers

func validateDateField(field FieldSpec) error {
	if field.Format == "" {
		return ourErrors.CNAB_ErrMissingDateFormat.Creator(fieldToError("Name", field.Name))
	}
	return nil
}

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

func formatDateField(field FieldSpec, value interface{}) (string, error) {
	dateValue, ok := value.(time.Time)
	if !ok {
		return "", ourErrors.CNAB_ErrFieldValueIsNotAnDate.Creator(fieldToError("Name", field.Name))
	}
	return formatDate(dateValue, field.Format), nil
}

// String Field Handlers

func validateStringField(field FieldSpec) error {
	// No specific validation needed for string fields
	return nil
}

func parseStringField(field FieldSpec, rawValue []byte) (interface{}, error) {
	return string(bytes.TrimSpace(rawValue)), nil
}

func formatStringField(field FieldSpec, value interface{}) (string, error) {
	strValue, ok := value.(string)
	if !ok {
		return "", ourErrors.CNAB_ErrFieldValueIsNotAnString.Creator(fieldToError("Name", field.Name))
	}
	return strValue, nil
}

// Utility Functions

var ErrInvalidIntegerInput = errors.New("invalid integer input")

func atoiUnsafe(b []byte) (int, error) {
	n := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, ErrInvalidIntegerInput
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func parseFloatBytes(b []byte, decimal int) (float64, error) {
	intValue, err := atoiUnsafe(b)
	if err != nil {
		return 0, err
	}
	return float64(intValue) / pow10(decimal), nil
}

var ErrInvalidDateLength = errors.New("invalid date length for field")

func parseDateBytes(b []byte, format string) (time.Time, error) {
	if len(b) != len(format) {
		return time.Time{}, ErrInvalidDateLength
	}
	goFormat := convertDateFormat(format)
	return time.Parse(goFormat, string(b))
}

func pow10(n int) float64 {
	result := 1.0
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

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

func formatDate(value time.Time, format string) string {
	goFormat := convertDateFormat(format)
	return value.Format(goFormat)
}

var (
	// ErrCannotConvertToInt is an error that occurs when a value cannot be converted to an int.
	ErrCannotConvertToInt   = ourErrors.CNAB_ErrCannotConvertToInt.Err
	IsErrCannotConvertToInt = iError.MatchError(ourErrors.CNAB_ErrCannotConvertToInt.Err)
)

func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, ourErrors.CNAB_ErrCannotConvertToInt.Creator(fieldToError("Value", fmt.Sprintf("%v", value)))
	}
}

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
		return strconv.ParseFloat(v, 64)
	default:
		return 0, ourErrors.CNAB_ErrCannotConvertToFloat.Creator(fieldToError("Value", fmt.Sprintf("%v", value)))
	}
}

func fieldToError(fieldName, fieldValue string) error {
	return fmt.Errorf("field %s %s", fieldName, fieldValue)
}
