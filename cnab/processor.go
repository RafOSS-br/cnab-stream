package cnab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// ParseRecord parses the provided record according to the loaded CNAB specification.
func (p *processor) ParseRecord(ctx context.Context, record []byte) (map[string]interface{}, error) {
	return nil, nil // Implement this
}

// PackRecord packs the provided data according to the loaded CNAB specification.
func (p *processor) PackRecord(ctx context.Context, data map[string]interface{}) ([]byte, error) {
	return nil, nil // Implement this
}
