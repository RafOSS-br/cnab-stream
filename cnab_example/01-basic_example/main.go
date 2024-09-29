package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RafOSS-br/cnab-stream/cnab"
)

func main() {

	// ------------
	// Setup
	// ------------

	ctx := context.Background()
	processor := cnab.NewProcessor()

	// ------------
	// Load spec
	// ------------

	specFile, err := os.Open("spec.json")
	if err != nil {
		log.Fatalf("Failed to open spec file: %v", err)
	}
	defer specFile.Close()

	err = processor.LoadSpec(ctx, specFile)
	if err != nil {
		log.Fatalf("Failed to load specification: %v", err)
	}

	// ------------
	// Print spec
	// ------------

	fmt.Println("Specification:")
	fmt.Println(processor.Spec())

	// ------------
	// Parse record
	// ------------

	// Corrected CNAB record with length 26
	record := []byte("0010001220240101121234")
	// 0010001202401010001234567 -> bank_code: 1, service_batch: 1, record_type: 2, payment_date: 2024-01-01, payment_amount: 12345.67

	fmt.Printf("Record length: %d\n", len(record)) // Should output 26

	data, err := processor.ParseRecord(ctx, record)
	if err != nil {
		log.Fatalf("Failed to parse record: %v", err)
	}

	fmt.Printf("Parsed Data: %+v\n", data)

	// Output:
	// Record length: 26
	// Parsed Data: map[bank_code:1 payment_amount:12345.67 payment_date:2024-01-01 00:00:00 +0000 UTC record_type:2 service_batch:1]

	// ------------
	// Pack record
	// ------------

	// Prepare data to pack
	dataToPack := map[string]interface{}{
		"bank_code":      1,
		"service_batch":  1,
		"record_type":    2,
		"payment_date":   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"payment_amount": 1235.67,
	}

	packedRecord, err := processor.PackRecord(ctx, dataToPack)
	if err != nil {
		log.Fatalf("Failed to pack record: %v", err)
	}

	fmt.Printf("Packed Record: %s\n", string(packedRecord))

	// Output:
	// Packed Record: 00100012202401011234567
}
