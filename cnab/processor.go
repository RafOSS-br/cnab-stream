package cnab

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	ourErrors "github.com/RafOSS-br/cnab-stream/errors"
	iError "github.com/RafOSS-br/cnab-stream/internal/error"
)

// Processor defines the interface for CNAB processing.
type Processor interface {
	LoadSpec(ctx context.Context, specReader io.Reader) error
	ParseRecord(ctx context.Context, record []byte) (map[string]interface{}, error)
	PackRecord(ctx context.Context, data map[string]interface{}) ([]byte, error)
}

type processor struct {
	spec       CNABSpec
	fieldCount int
}

// NewProcessor creates a new CNAB processor.
func NewProcessor() Processor {
	return &processor{}
}

// FieldSpec defines the specification for a single field.
type FieldSpec struct {
	Name    string `json:"name"`
	Type    string `json:"type"` // "int", "float", "date", "string"
	Start   int    `json:"start"`
	Length  int    `json:"length"`
	Format  string `json:"format,omitempty"`
	Decimal int    `json:"decimal,omitempty"`
	End     int    // Calculated as Start + Length - 1
}

// CNABSpec defines the CNAB specification.
type CNABSpec struct {
	Fields []FieldSpec `json:"fields"`
}

var (
	// ErrFailedToDecodeSpecJSON is an error that occurs when the CNAB spec JSON cannot be decoded.
	ErrFailedToDecodeSpecJSON   = ourErrors.CNAB_ErrFailedToDecodeSpecJSON.Err
	IsErrFailedToDecodeSpecJSON = iError.MatchError(ourErrors.CNAB_ErrFailedToDecodeSpecJSON.Err)
	// ErrStartAndLengthMustBeGreaterThanZero is an error that occurs when the start and length of a field are less than or equal to zero.
	ErrStartAndLengthMustBeGreaterThanZero   = ourErrors.CNAB_ErrStartAndLengthMustBeGreaterThanZeroEncapsulator.Err
	IsErrStartAndLengthMustBeGreaterThanZero = iError.MatchError(ourErrors.CNAB_ErrStartAndLengthMustBeGreaterThanZeroEncapsulator.Err)
	// ErrFieldHasNoTypeSpecified is an error that occurs when a field has no type specified.
	ErrFieldHasNoTypeSpecified   = ourErrors.CNAB_ErrFieldHasNoTypeSpecified.Err
	IsErrFieldHasNoTypeSpecified = iError.MatchError(ourErrors.CNAB_ErrFieldHasNoTypeSpecified.Err)
)

// LoadSpec loads the CNAB specification from a JSON reader.
func (p *processor) LoadSpec(ctx context.Context, specReader io.Reader) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	decoder := json.NewDecoder(specReader)
	if err := decoder.Decode(&p.spec); err != nil {
		return ourErrors.CNAB_ErrFailedToDecodeSpecJSON.Creator(err)
	}

	// Precompute field positions and validate fields
	for i := range p.spec.Fields {
		field := &p.spec.Fields[i]
		field.End = field.Start + field.Length - 1

		if field.Start <= 0 || field.Length <= 0 {
			return ourErrors.CNAB_ErrStartAndLengthMustBeGreaterThanZeroEncapsulator.Creator(fieldToError("Name", field.Name))
		}
		if field.Type == "" {
			return ourErrors.CNAB_ErrFieldHasNoTypeSpecified.Creator(fieldToError("Name", field.Name))
		}
	}

	p.fieldCount = len(p.spec.Fields)
	return nil
}

// ParseRecord parses a CNAB record into a map.
func (p *processor) ParseRecord(ctx context.Context, record []byte) (map[string]interface{}, error) {
	result := make(map[string]interface{}, p.fieldCount)

	for _, field := range p.spec.Fields {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if err := p.parseRecord(record, field, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Record Pool
var recordPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 1024) // Adjust capacity as needed
	},
}

var (
	// ErrMissingDataForField is an error that occurs when data is missing for a field.
	ErrMissingDataForField   = ourErrors.CNAB_ErrMissingDataForField.Err
	IsErrMissingDataForField = iError.MatchError(ourErrors.CNAB_ErrMissingDataForField.Err)

	// ErrFailedToFormatField is an error that occurs when a field cannot be formatted.
	ErrFailedToFormatField   = ourErrors.CNAB_ErrFailedToFormatField.Err
	IsErrFailedToFormatField = iError.MatchError(ourErrors.CNAB_ErrFailedToFormatField.Err)

	// ErrFieldExceedsSpecifiedLength is an error that occurs when a field exceeds the specified length.
	ErrFieldExceedsSpecifiedLength   = ourErrors.CNAB_ErrFieldExceedsSpecifiedLength.Err
	IsErrFieldExceedsSpecifiedLength = iError.MatchError(ourErrors.CNAB_ErrFieldExceedsSpecifiedLength.Err)
)

// PackRecord packs data into a CNAB record.
func (p *processor) PackRecord(ctx context.Context, data map[string]interface{}) ([]byte, error) {
	// Calculate total length
	totalLength := 0
	for _, field := range p.spec.Fields {
		totalLength += field.Length
	}

	// Get buffer from pool
	buf := recordPool.Get().([]byte)
	defer func() {
		buf = buf[:0]
		//lint:ignore SA6002 recordPool.Put is safe to call with buf. Only header is copied, slice store reference to data
		recordPool.Put(buf)
	}()

	if cap(buf) < totalLength {
		buf = make([]byte, totalLength)
	} else {
		buf = buf[:totalLength]
	}

	for _, field := range p.spec.Fields {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		value, exists := data[field.Name]
		if !exists {
			return nil, ourErrors.CNAB_ErrMissingDataForField.Creator(fieldToError("Name", field.Name))
		}

		strValue, err := p.formatFieldValue(field, value)
		if err != nil {
			return nil, ourErrors.CNAB_ErrFailedToFormatField.Creator(err)
		}

		if len(strValue) > field.Length {
			return nil, ourErrors.CNAB_ErrFieldExceedsSpecifiedLength.Creator(fieldToError("Name", field.Name))
		}

		// Fill with spaces if necessary
		paddedValue := strValue + strings.Repeat(" ", field.Length-len(strValue))
		copy(buf[field.Start-1:field.End], paddedValue)
	}

	return buf, nil
}

// Helper Methods

var (
	// ErrFieldIsEmpty is an error that occurs when a field is empty.
	ErrFieldIsEmpty   = ourErrors.CNAB_ErrFieldIsEmpty.Err
	IsErrFieldIsEmpty = iError.MatchError(ourErrors.CNAB_ErrFieldIsEmpty.Err)

	// ErrUnsupportedFieldType is an error that occurs when a field has an unsupported type.
	ErrUnsupportedFieldType   = ourErrors.CNAB_ErrUnsupportedFieldType.Err
	IsErrUnsupportedFieldType = iError.MatchError(ourErrors.CNAB_ErrUnsupportedFieldType.Err)
)

func (p *processor) parseFieldValue(field FieldSpec, rawValue []byte) (interface{}, error) {
	trimmedValue := bytes.TrimSpace(rawValue)

	// Basic input sanitization: ensure trimmedValue is not empty
	if len(trimmedValue) == 0 {
		return nil, ourErrors.CNAB_ErrFieldIsEmpty.Creator(fieldToError("Name", field.Name))
	}

	switch field.Type {
	case "int":
		return atoiUnsafe(trimmedValue)
	case "float":
		return parseFloatBytes(trimmedValue, field.Decimal)
	case "date":
		return parseDateBytes(trimmedValue, field.Format)
	case "string":
		return string(trimmedValue), nil
	default:
		return nil, ourErrors.CNAB_ErrUnsupportedFieldType.Creator(fieldToError("Type", field.Type))
	}
}

var (
	// ErrFieldValueIsNotAnDate is an error that occurs when a field value is not a date.
	ErrFieldValueIsNotAnDate   = ourErrors.CNAB_ErrFieldValueIsNotAnDate.Err
	IsErrFieldValueIsNotAnDate = iError.MatchError(ourErrors.CNAB_ErrFieldValueIsNotAnDate.Err)

	// ErrFieldValueIsNotAnString is an error that occurs when a field value is not a string.
	ErrFieldValueIsNotAnString   = ourErrors.CNAB_ErrFieldValueIsNotAnString.Err
	IsErrFieldValueIsNotAnString = iError.MatchError(ourErrors.CNAB_ErrFieldValueIsNotAnString.Err)

	// ErrFieldValueIsNotAnInt is an error that occurs when a field value is not an int.
	ErrFieldValueIsNotAnInt   = ourErrors.CNAB_ErrFieldValueIsNotAnInt.Err
	IsErrFieldValueIsNotAnInt = iError.MatchError(ourErrors.CNAB_ErrFieldValueIsNotAnInt.Err)

	// ErrFieldValueIsNotAnFloat is an error that occurs when a field value is not a float.
	ErrFieldValueIsNotAnFloat   = ourErrors.CNAB_ErrFieldValueIsNotAnFloat.Err
	IsErrFieldValueIsNotAnFloat = iError.MatchError(ourErrors.CNAB_ErrFieldValueIsNotAnFloat.Err)
)

func (p *processor) formatFieldValue(field FieldSpec, value interface{}) (string, error) {
	switch field.Type {
	case "int":
		intValue, err := toInt(value)
		if err != nil {
			return "", ourErrors.CNAB_ErrFieldValueIsNotAnInt.Creator(fmt.Errorf("field %s: %w", field.Name, err))
		}
		return fmt.Sprintf("%0*d", field.Length, intValue), nil
	case "float":
		floatValue, err := toFloat(value)
		if err != nil {
			return "", ourErrors.CNAB_ErrFieldValueIsNotAnFloat.Creator(fieldToError("Name", field.Name))
		}
		scaledValue := int64(floatValue * pow10(field.Decimal))
		return fmt.Sprintf("%0*d", field.Length, scaledValue), nil
	case "date":
		dateValue, ok := value.(time.Time)
		if !ok {
			return "", ourErrors.CNAB_ErrFieldValueIsNotAnDate.Creator(fieldToError("Name", field.Name))
		}
		return formatDate(dateValue, field.Format), nil
	case "string":
		strValue, ok := value.(string)
		if !ok {
			return "", ourErrors.CNAB_ErrFieldValueIsNotAnString.Creator(fieldToError("Name", field.Name))
		}
		return strValue, nil
	default:
		return "", ourErrors.CNAB_ErrUnsupportedFieldType.Creator(fieldToError("Type", field.Type))
	}
}

var (
	// ErrFieldExceedsRecordLength is an error that occurs when a field exceeds the record length.
	ErrFieldExceedsRecordLength   = ourErrors.CNAB_ErrFieldExceedsRecordLength.Err
	IsErrFieldExceedsRecordLength = iError.MatchError(ourErrors.CNAB_ErrFieldExceedsRecordLength.Err)

	// ErrFailedToParseField is an error that occurs when a field cannot be parsed.
	ErrFailedToParseField   = ourErrors.CNAB_ErrFailedToParseField.Err
	IsErrFailedToParseField = iError.MatchError(ourErrors.CNAB_ErrFailedToParseField.Err)
)

// parse a CNAB record into a map
func (p *processor) parseRecord(record []byte, field FieldSpec, m map[string]interface{}) error {
	if field.End > len(record) {
		return ourErrors.CNAB_ErrFieldExceedsRecordLength.Creator(fieldToError("Name", field.Name))
	}

	rawValue := record[field.Start-1 : field.End]
	value, err := p.parseFieldValue(field, rawValue)
	if err != nil {
		return ourErrors.CNAB_ErrFailedToParseField.Creator(fmt.Errorf("field %s: %w", field.Name, err))
	}

	m[field.Name] = value
	return nil
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
