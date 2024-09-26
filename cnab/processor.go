package cnab

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
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

// LoadSpec loads the CNAB specification from a JSON reader.
func (p *processor) LoadSpec(ctx context.Context, specReader io.Reader) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	decoder := json.NewDecoder(specReader)
	if err := decoder.Decode(&p.spec); err != nil {
		return fmt.Errorf("failed to decode spec JSON: %w", err)
	}

	// Precompute field positions and validate fields
	for i := range p.spec.Fields {
		field := &p.spec.Fields[i]
		field.End = field.Start + field.Length - 1

		if field.Start <= 0 || field.Length <= 0 {
			return fmt.Errorf("invalid field specification for %s: start and length must be positive", field.Name)
		}
		if field.Type == "" {
			return fmt.Errorf("field %s has no type specified", field.Name)
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

		if field.End > len(record) {
			return nil, fmt.Errorf("field %s exceeds record length", field.Name)
		}

		rawValue := record[field.Start-1 : field.End]
		value, err := p.parseFieldValue(field, rawValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field %s: %w", field.Name, err)
		}

		result[field.Name] = value
	}

	return result, nil
}

// PackRecord packs the provided data according to the loaded CNAB specification.
func (p *processor) PackRecord(ctx context.Context, data map[string]interface{}) ([]byte, error) {
	return nil, nil // Implement this
}

// Helper Methods

func (p *processor) parseFieldValue(field FieldSpec, rawValue []byte) (interface{}, error) {
	trimmedValue := bytes.TrimSpace(rawValue)

	// Basic input sanitization: ensure trimmedValue is not empty
	if len(trimmedValue) == 0 {
		return nil, fmt.Errorf("field %s is empty", field.Name)
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
		return nil, fmt.Errorf("unsupported field type: %s", field.Type)
	}
}

// Utility Functions

func atoiUnsafe(b []byte) (int, error) {
	n := 0
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, errors.New("invalid integer input")
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

func parseDateBytes(b []byte, format string) (time.Time, error) {
	if len(b) != len(format) {
		return time.Time{}, fmt.Errorf("invalid date length for field")
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
