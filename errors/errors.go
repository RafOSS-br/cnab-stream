package errors

import (
	iError "github.com/RafOSS-br/cnab-stream/internal/error"
)

type ErrInstance[T any] struct {
	Encapsulator iError.Encapsulator[T]
	Err          error
	Creator      func(error) error
}

type FailedToDecodeSpecJSONEncapsulator struct {
	iError.InternalError
}

func (e *FailedToDecodeSpecJSONEncapsulator) CreateError(err error) error {
	return &FailedToDecodeSpecJSONEncapsulator{iError.NewError("Failed to decode spec JSON", err)}
}

type StartAndLengthMustBeGreaterThanZero struct {
	iError.InternalError
}

func (e *StartAndLengthMustBeGreaterThanZero) CreateError(err error) error {
	return &StartAndLengthMustBeGreaterThanZero{iError.NewError("Start and length must be greater than zero", err)}
}

type FieldHasNoTypeSpecified struct {
	iError.InternalError
}

func (e *FieldHasNoTypeSpecified) CreateError(err error) error {
	return &FieldHasNoTypeSpecified{iError.NewError("Field has no type specified", err)}
}

type MissingDataForField struct {
	iError.InternalError
}

func (e *MissingDataForField) CreateError(err error) error {
	return &MissingDataForField{iError.NewError("Missing data for field", err)}
}

type FailedToFormatField struct {
	iError.InternalError
}

func (e *FailedToFormatField) CreateError(err error) error {
	return &FailedToFormatField{iError.NewError("Failed to format field", err)}
}

type FieldExceedsSpecifiedLength struct {
	iError.InternalError
}

func (e *FieldExceedsSpecifiedLength) CreateError(err error) error {
	return &FieldExceedsSpecifiedLength{iError.NewError("Field exceeds specified length", err)}
}

type FieldIsEmpty struct {
	iError.InternalError
}

func (e *FieldIsEmpty) CreateError(err error) error {
	return &FieldIsEmpty{iError.NewError("Field is empty", err)}
}

type UnsupportedFieldType struct {
	iError.InternalError
}

func (e *UnsupportedFieldType) CreateError(err error) error {
	return &UnsupportedFieldType{iError.NewError("Unsupported field type", err)}
}

type FieldValueIsNotAnDate struct {
	iError.InternalError
}

func (e *FieldValueIsNotAnDate) CreateError(err error) error {
	return &FieldValueIsNotAnDate{iError.NewError("Field value is not an date", err)}
}

type FieldValueIsNotAnString struct {
	iError.InternalError
}

func (e *FieldValueIsNotAnString) CreateError(err error) error {
	return &FieldValueIsNotAnString{iError.NewError("Field value is not an string", err)}
}

type FieldValueIsNotAnInt struct {
	iError.InternalError
}

func (e *FieldValueIsNotAnInt) CreateError(err error) error {
	return &FieldValueIsNotAnInt{iError.NewError("Field value is not an int", err)}
}

type FieldValueIsNotAnFloat struct {
	iError.InternalError
}

func (e *FieldValueIsNotAnFloat) CreateError(err error) error {
	return &FieldValueIsNotAnFloat{iError.NewError("Field value is not an float", err)}
}

type FieldExceedsRecordLength struct {
	iError.InternalError
}

func (e *FieldExceedsRecordLength) CreateError(err error) error {
	return &FieldExceedsRecordLength{iError.NewError("Field exceeds record length", err)}
}

type FailedToParseField struct {
	iError.InternalError
}

func (e *FailedToParseField) CreateError(err error) error {
	return &FailedToParseField{iError.NewError("Failed to parse field", err)}
}

type CannotConvertToInt struct {
	iError.InternalError
}

func (e *CannotConvertToInt) CreateError(err error) error {
	return &CannotConvertToInt{iError.NewError("Cannot convert to int", err)}
}

type CannotConvertToFloat struct {
	iError.InternalError
}

func (e *CannotConvertToFloat) CreateError(err error) error {
	return &CannotConvertToFloat{iError.NewError("Cannot convert to float", err)}
}

type InvalidDecimalValue struct {
	iError.InternalError
}

func (e *InvalidDecimalValue) CreateError(err error) error {
	return &InvalidDecimalValue{iError.NewError("Invalid decimal value", err)}
}

type MissingDateFormat struct {
	iError.InternalError
}

func (e *MissingDateFormat) CreateError(err error) error {
	return &MissingDateFormat{iError.NewError("Missing date format", err)}
}

type InvalidDateLength struct {
	iError.InternalError
}

func (e *InvalidDateLength) CreateError(err error) error {
	return &InvalidDateLength{iError.NewError("Invalid date length", err)}
}

type StartMustBeGreaterOrEqualZero struct {
	iError.InternalError
}

func (e *StartMustBeGreaterOrEqualZero) CreateError(err error) error {
	return &StartMustBeGreaterOrEqualZero{iError.NewError("Start must be greater or equal zero", err)}
}

type CancelledContext struct {
	iError.InternalError
}

func (e *CancelledContext) CreateError(err error) error {
	return &CancelledContext{iError.NewError("Canceled context", err)}
}

type FormatError struct {
	iError.InternalError
}

func (e *FormatError) CreateError(err error) error {
	return &FormatError{iError.NewError("Format error", err)}
}

var (
	// CBAN errors

	// CNAB_ErrFailedToDecodeSpecJSON is an error that occurs when the CNAB spec JSON cannot be decoded.
	CNAB_ErrFailedToDecodeSpecJSON = &ErrInstance[*FailedToDecodeSpecJSONEncapsulator]{
		Encapsulator: &FailedToDecodeSpecJSONEncapsulator{},
		Creator:      iError.NewCreator[*FailedToDecodeSpecJSONEncapsulator](&FailedToDecodeSpecJSONEncapsulator{}),
		Err:          getFirstClass(iError.NewCreator[*FailedToDecodeSpecJSONEncapsulator](&FailedToDecodeSpecJSONEncapsulator{})),
	}

	// CNAB_ErrLengthMustBeGreaterThanZero is an error that occurs when the start and length of a field are less than or equal to zero.
	CNAB_ErrLengthMustBeGreaterThanZeroEncapsulator = &ErrInstance[*StartAndLengthMustBeGreaterThanZero]{
		Encapsulator: &StartAndLengthMustBeGreaterThanZero{},
		Creator:      iError.NewCreator[*StartAndLengthMustBeGreaterThanZero](&StartAndLengthMustBeGreaterThanZero{}),
		Err:          getFirstClass(iError.NewCreator[*StartAndLengthMustBeGreaterThanZero](&StartAndLengthMustBeGreaterThanZero{})),
	}

	// CNAB_ErrFieldHasNoTypeSpecified is an error that occurs when a field has no type specified.
	CNAB_ErrFieldHasNoTypeSpecified = &ErrInstance[*FieldHasNoTypeSpecified]{
		Encapsulator: &FieldHasNoTypeSpecified{},
		Creator:      iError.NewCreator[*FieldHasNoTypeSpecified](&FieldHasNoTypeSpecified{}),
		Err:          getFirstClass(iError.NewCreator[*FieldHasNoTypeSpecified](&FieldHasNoTypeSpecified{})),
	}

	// CNAB_ErrMissingDataForField is an error that occurs when a field is missing data.
	CNAB_ErrMissingDataForField = &ErrInstance[*MissingDataForField]{
		Encapsulator: &MissingDataForField{},
		Creator:      iError.NewCreator[*MissingDataForField](&MissingDataForField{}),
		Err:          getFirstClass(iError.NewCreator[*MissingDataForField](&MissingDataForField{})),
	}

	// CNAB_ErrFailedToFormatField is an error that occurs when a field cannot be formatted.
	CNAB_ErrFailedToFormatField = &ErrInstance[*FailedToFormatField]{
		Encapsulator: &FailedToFormatField{},
		Creator:      iError.NewCreator[*FailedToFormatField](&FailedToFormatField{}),
		Err:          getFirstClass(iError.NewCreator[*FailedToFormatField](&FailedToFormatField{})),
	}

	// CNAB_ErrFieldExceedsSpecifiedLength is an error that occurs when a field exceeds the specified length.
	CNAB_ErrFieldExceedsSpecifiedLength = &ErrInstance[*FieldExceedsSpecifiedLength]{
		Encapsulator: &FieldExceedsSpecifiedLength{},
		Creator:      iError.NewCreator[*FieldExceedsSpecifiedLength](&FieldExceedsSpecifiedLength{}),
		Err:          getFirstClass(iError.NewCreator[*FieldExceedsSpecifiedLength](&FieldExceedsSpecifiedLength{})),
	}

	// CNAB_ErrFieldIsEmpty is an error that occurs when a field is empty.
	CNAB_ErrFieldIsEmpty = &ErrInstance[*FieldIsEmpty]{
		Encapsulator: &FieldIsEmpty{},
		Creator:      iError.NewCreator[*FieldIsEmpty](&FieldIsEmpty{}),
		Err:          getFirstClass(iError.NewCreator[*FieldIsEmpty](&FieldIsEmpty{})),
	}

	// CNAB_ErrUnsupportedFieldType is an error that occurs when a field type is unsupported.
	CNAB_ErrUnsupportedFieldType = &ErrInstance[*UnsupportedFieldType]{
		Encapsulator: &UnsupportedFieldType{},
		Creator:      iError.NewCreator[*UnsupportedFieldType](&UnsupportedFieldType{}),
		Err:          getFirstClass(iError.NewCreator[*UnsupportedFieldType](&UnsupportedFieldType{})),
	}

	// CNAB_ErrFieldValueIsNotAnDate is an error that occurs when a field value is not a date.
	CNAB_ErrFieldValueIsNotAnDate = &ErrInstance[*FieldValueIsNotAnDate]{
		Encapsulator: &FieldValueIsNotAnDate{},
		Creator:      iError.NewCreator[*FieldValueIsNotAnDate](&FieldValueIsNotAnDate{}),
		Err:          getFirstClass(iError.NewCreator[*FieldValueIsNotAnDate](&FieldValueIsNotAnDate{})),
	}

	// CNAB_ErrFieldValueIsNotAnString is an error that occurs when a field value is not a string.
	CNAB_ErrFieldValueIsNotAnString = &ErrInstance[*FieldValueIsNotAnString]{
		Encapsulator: &FieldValueIsNotAnString{},
		Creator:      iError.NewCreator[*FieldValueIsNotAnString](&FieldValueIsNotAnString{}),
		Err:          getFirstClass(iError.NewCreator[*FieldValueIsNotAnString](&FieldValueIsNotAnString{})),
	}

	// CNAB_ErrFieldValueIsNotAnInt is an error that occurs when a field value is not an int.
	CNAB_ErrFieldValueIsNotAnInt = &ErrInstance[*FieldValueIsNotAnInt]{
		Encapsulator: &FieldValueIsNotAnInt{},
		Creator:      iError.NewCreator[*FieldValueIsNotAnInt](&FieldValueIsNotAnInt{}),
		Err:          getFirstClass(iError.NewCreator[*FieldValueIsNotAnInt](&FieldValueIsNotAnInt{})),
	}

	// CNAB_ErrFieldValueIsNotAnFloat is an error that occurs when a field value is not a float.
	CNAB_ErrFieldValueIsNotAnFloat = &ErrInstance[*FieldValueIsNotAnFloat]{
		Encapsulator: &FieldValueIsNotAnFloat{},
		Creator:      iError.NewCreator[*FieldValueIsNotAnFloat](&FieldValueIsNotAnFloat{}),
		Err:          getFirstClass(iError.NewCreator[*FieldValueIsNotAnFloat](&FieldValueIsNotAnFloat{})),
	}

	// CNAB_ErrFieldExceedsRecordLength is an error that occurs when a field exceeds the record length.
	CNAB_ErrFieldExceedsRecordLength = &ErrInstance[*FieldExceedsRecordLength]{
		Encapsulator: &FieldExceedsRecordLength{},
		Creator:      iError.NewCreator[*FieldExceedsRecordLength](&FieldExceedsRecordLength{}),
		Err:          getFirstClass(iError.NewCreator[*FieldExceedsRecordLength](&FieldExceedsRecordLength{})),
	}

	// CNAB_ErrFailedToParseField is an error that occurs when a field cannot be parsed.
	CNAB_ErrFailedToParseField = &ErrInstance[*FailedToParseField]{
		Encapsulator: &FailedToParseField{},
		Creator:      iError.NewCreator[*FailedToParseField](&FailedToParseField{}),
		Err:          getFirstClass(iError.NewCreator[*FailedToParseField](&FailedToParseField{})),
	}

	// CNAB_ErrInvalidDecimalValue is an error that occurs when a decimal value is invalid.
	CNAB_ErrInvalidDecimalValue = &ErrInstance[*InvalidDecimalValue]{
		Encapsulator: &InvalidDecimalValue{},
		Creator:      iError.NewCreator[*InvalidDecimalValue](&InvalidDecimalValue{}),
		Err:          getFirstClass(iError.NewCreator[*InvalidDecimalValue](&InvalidDecimalValue{})),
	}

	// CNAB_ErrMissingDateFormat is an error that occurs when a date format is missing.
	CNAB_ErrMissingDateFormat = &ErrInstance[*MissingDateFormat]{
		Encapsulator: &MissingDateFormat{},
		Creator:      iError.NewCreator[*MissingDateFormat](&MissingDateFormat{}),
		Err:          getFirstClass(iError.NewCreator[*MissingDateFormat](&MissingDateFormat{})),
	}

	// CNAB_ErrInvalidDateLength is an error that occurs when a date length is invalid.
	CNAB_ErrInvalidDateLength = &ErrInstance[*InvalidDateLength]{
		Encapsulator: &InvalidDateLength{},
		Creator:      iError.NewCreator[*InvalidDateLength](&InvalidDateLength{}),
		Err:          getFirstClass(iError.NewCreator[*InvalidDateLength](&InvalidDateLength{})),
	}

	// CNAB_ErrCancelledContext is an error that occurs when a context is canceled.
	CNAB_ErrCancelledContext = &ErrInstance[*CancelledContext]{
		Encapsulator: &CancelledContext{},
		Creator:      iError.NewCreator[*CancelledContext](&CancelledContext{}),
		Err:          getFirstClass(iError.NewCreator[*CancelledContext](&CancelledContext{})),
	}

	// CNAB_ErrStartMustBeGreaterOrEqualZero is an error that occurs when the start is less than zero.
	CNAB_ErrStartMustBeGreaterOrEqualZero = &ErrInstance[*StartMustBeGreaterOrEqualZero]{
		Encapsulator: &StartMustBeGreaterOrEqualZero{},
		Creator:      iError.NewCreator[*StartMustBeGreaterOrEqualZero](&StartMustBeGreaterOrEqualZero{}),
		Err:          getFirstClass(iError.NewCreator[*StartMustBeGreaterOrEqualZero](&StartMustBeGreaterOrEqualZero{})),
	}

	// CNAB_ErrFormatError is an error that occurs when a format error occurs.
	CNAB_ErrFormatError = &ErrInstance[*FormatError]{
		Encapsulator: &FormatError{},
		Creator:      iError.NewCreator[*FormatError](&FormatError{}),
		Err:          getFirstClass(iError.NewCreator[*FormatError](&FormatError{})),
	}
)

// export helper

func getFirstClass(f func(error) error) error {
	err := f(nil)
	if iErr, ok := err.(iError.InternalError); ok {
		return iErr.FirstClass()
	}
	return err
}
