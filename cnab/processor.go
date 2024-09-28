package cnab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	ourErrors "github.com/RafOSS-br/cnab-stream/errors"
	iError "github.com/RafOSS-br/cnab-stream/internal/error"
)

// Processor defines the interface for CNAB processing.
type Processor interface {
	LoadSpec(ctx context.Context, specReader io.Reader) error
	ParseRecord(ctx context.Context, record []byte) (map[string]interface{}, error)
	PackRecord(ctx context.Context, data map[string]interface{}) ([]byte, error)
}

// processor implements the Processor interface.
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

// FieldHandler defines the interface for handling fields.
type FieldHandler struct {
	Validate func(field FieldSpec) error
	Parse    func(field FieldSpec, rawValue []byte) (interface{}, error)
	Format   func(field FieldSpec, value interface{}) (string, error)
}

// Field Handlers
var fieldHandlers = map[string]*FieldHandler{
	"int": {
		Validate: validateIntField,
		Parse:    parseIntField,
		Format:   formatIntField,
	},
	"float": {
		Validate: validateFloatField,
		Parse:    parseFloatField,
		Format:   formatFloatField,
	},
	"date": {
		Validate: validateDateField,
		Parse:    parseDateField,
		Format:   formatDateField,
	},
	"string": {
		Validate: validateStringField,
		Parse:    parseStringField,
		Format:   formatStringField,
	},
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

	// ErrCancelledContext is an error that occurs when the context is cancelled.
	ErrCancelledContext   = ourErrors.CNAB_ErrCancelledContext.Err
	IsErrCancelledContext = iError.MatchError(ourErrors.CNAB_ErrCancelledContext.Err)
)

// LoadSpec loads the CNAB specification from a JSON reader.
func (p *processor) LoadSpec(ctx context.Context, specReader io.Reader) error {
	select {
	case <-ctx.Done():
		return ourErrors.CNAB_ErrCancelledContext.Creator(ctx.Err())
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
		if _, exists := fieldHandlers[field.Type]; !exists {
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
			return nil, ourErrors.CNAB_ErrCancelledContext.Creator(ctx.Err())
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
			return nil, ourErrors.CNAB_ErrCancelledContext.Creator(ctx.Err())
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

// formatFieldValue formats a field value.
func (p *processor) formatFieldValue(field FieldSpec, value interface{}) (string, error) {
	handler := fieldHandlers[field.Type]
	if handler == nil {
		return "", ourErrors.CNAB_ErrUnsupportedFieldType.Creator(fieldToError("Type", field.Type))
	}

	return handler.Format(field, value)
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
	handler := fieldHandlers[field.Type]
	if handler == nil {
		return ourErrors.CNAB_ErrUnsupportedFieldType.Creator(fieldToError("Type", field.Type))
	}

	if field.End > len(record) {
		return ourErrors.CNAB_ErrFieldExceedsRecordLength.Creator(fieldToError("Name", field.Name))
	}

	rawValue := record[field.Start-1 : field.End]
	value, err := handler.Parse(field, rawValue)
	if err != nil {
		return ourErrors.CNAB_ErrFailedToParseField.Creator(fmt.Errorf("field %s: %w", field.Name, err))
	}

	m[field.Name] = value
	return nil
}
