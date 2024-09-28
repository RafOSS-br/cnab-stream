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
	*iError.IError
}

func (e *FailedToDecodeSpecJSONEncapsulator) CreateError(err error) error {
	return &FailedToDecodeSpecJSONEncapsulator{iError.NewError("Failed to decode spec JSON", err)}
}

type StartAndLengthMustBeGreaterThanZero struct {
	*iError.IError
}

func (e *StartAndLengthMustBeGreaterThanZero) CreateError(err error) error {
	return &StartAndLengthMustBeGreaterThanZero{iError.NewError("Start and length must be greater than zero", err)}
}

type FieldHasNoTypeSpecified struct {
	*iError.IError
}

func (e *FieldHasNoTypeSpecified) CreateError(err error) error {
	return &FieldHasNoTypeSpecified{iError.NewError("Field has no type specified", err)}
}

type MissingDataForField struct {
	*iError.IError
}

func (e *MissingDataForField) CreateError(err error) error {
	return &MissingDataForField{iError.NewError("Missing data for field", err)}
}

type FailedToFormatField struct {
	*iError.IError
}

func (e *FailedToFormatField) CreateError(err error) error {
	return &FailedToFormatField{iError.NewError("Failed to format field", err)}
}

type FieldExceedsSpecifiedLength struct {
	*iError.IError
}

func (e *FieldExceedsSpecifiedLength) CreateError(err error) error {
	return &FieldExceedsSpecifiedLength{iError.NewError("Field exceeds specified length", err)}
}

type FieldIsEmpty struct {
	*iError.IError
}

func (e *FieldIsEmpty) CreateError(err error) error {
	return &FieldIsEmpty{iError.NewError("Field is empty", err)}
}

type UnsupportedFieldType struct {
	*iError.IError
}

func (e *UnsupportedFieldType) CreateError(err error) error {
	return &UnsupportedFieldType{iError.NewError("Unsupported field type", err)}
}

type FieldValueIsNotAnDate struct {
	*iError.IError
}

func (e *FieldValueIsNotAnDate) CreateError(err error) error {
	return &FieldValueIsNotAnDate{iError.NewError("Field value is not an date", err)}
}

type FieldValueIsNotAnString struct {
	*iError.IError
}

func (e *FieldValueIsNotAnString) CreateError(err error) error {
	return &FieldValueIsNotAnString{iError.NewError("Field value is not an string", err)}
}

type FieldValueIsNotAnInt struct {
	*iError.IError
}

func (e *FieldValueIsNotAnInt) CreateError(err error) error {
	return &FieldValueIsNotAnInt{iError.NewError("Field value is not an int", err)}
}

type FieldValueIsNotAnFloat struct {
	*iError.IError
}

func (e *FieldValueIsNotAnFloat) CreateError(err error) error {
	return &FieldValueIsNotAnFloat{iError.NewError("Field value is not an float", err)}
}

type FieldExceedsRecordLength struct {
	*iError.IError
}

func (e *FieldExceedsRecordLength) CreateError(err error) error {
	return &FieldExceedsRecordLength{iError.NewError("Field exceeds record length", err)}
}

type FailedToParseField struct {
	*iError.IError
}

func (e *FailedToParseField) CreateError(err error) error {
	return &FailedToParseField{iError.NewError("Failed to parse field", err)}
}

type CannotConvertToInt struct {
	*iError.IError
}

func (e *CannotConvertToInt) CreateError(err error) error {
	return &CannotConvertToInt{iError.NewError("Cannot convert to int", err)}
}

type CannotConvertToFloat struct {
	*iError.IError
}

func (e *CannotConvertToFloat) CreateError(err error) error {
	return &CannotConvertToFloat{iError.NewError("Cannot convert to float", err)}
}

type InvalidDecimalValue struct {
	*iError.IError
}

func (e *InvalidDecimalValue) CreateError(err error) error {
	return &InvalidDecimalValue{iError.NewError("Invalid decimal value", err)}
}

type MissingDateFormat struct {
	*iError.IError
}

func (e *MissingDateFormat) CreateError(err error) error {
	return &MissingDateFormat{iError.NewError("Missing date format", err)}
}

type InvalidDateLength struct {
	*iError.IError
}

func (e *InvalidDateLength) CreateError(err error) error {
	return &InvalidDateLength{iError.NewError("Invalid date length", err)}
}

var (
	// CBAN errors

	// CNAB_ErrFailedToDecodeSpecJSON is an error that occurs when the CNAB spec JSON cannot be decoded.
	CNAB_ErrFailedToDecodeSpecJSON = &ErrInstance[*FailedToDecodeSpecJSONEncapsulator]{
		Encapsulator: &FailedToDecodeSpecJSONEncapsulator{},
		Creator:      iError.NewCreator[*FailedToDecodeSpecJSONEncapsulator](&FailedToDecodeSpecJSONEncapsulator{}),
		Err:          iError.NewCreator[*FailedToDecodeSpecJSONEncapsulator](&FailedToDecodeSpecJSONEncapsulator{})(nil),
	}

	// CNAB_ErrStartAndLengthMustBeGreaterThanZero is an error that occurs when the start and length of a field are less than or equal to zero.
	CNAB_ErrStartAndLengthMustBeGreaterThanZeroEncapsulator = &ErrInstance[*StartAndLengthMustBeGreaterThanZero]{
		Encapsulator: &StartAndLengthMustBeGreaterThanZero{},
		Creator:      iError.NewCreator[*StartAndLengthMustBeGreaterThanZero](&StartAndLengthMustBeGreaterThanZero{}),
		Err:          iError.NewCreator[*StartAndLengthMustBeGreaterThanZero](&StartAndLengthMustBeGreaterThanZero{})(nil),
	}

	// CNAB_ErrFieldHasNoTypeSpecified is an error that occurs when a field has no type specified.
	CNAB_ErrFieldHasNoTypeSpecified = &ErrInstance[*FieldHasNoTypeSpecified]{
		Encapsulator: &FieldHasNoTypeSpecified{},
		Creator:      iError.NewCreator[*FieldHasNoTypeSpecified](&FieldHasNoTypeSpecified{}),
		Err:          iError.NewCreator[*FieldHasNoTypeSpecified](&FieldHasNoTypeSpecified{})(nil),
	}

	// CNAB_ErrMissingDataForField is an error that occurs when a field is missing data.
	CNAB_ErrMissingDataForField = &ErrInstance[*MissingDataForField]{
		Encapsulator: &MissingDataForField{},
		Creator:      iError.NewCreator[*MissingDataForField](&MissingDataForField{}),
		Err:          iError.NewCreator[*MissingDataForField](&MissingDataForField{})(nil),
	}

	// CNAB_ErrFailedToFormatField is an error that occurs when a field cannot be formatted.
	CNAB_ErrFailedToFormatField = &ErrInstance[*FailedToFormatField]{
		Encapsulator: &FailedToFormatField{},
		Creator:      iError.NewCreator[*FailedToFormatField](&FailedToFormatField{}),
		Err:          iError.NewCreator[*FailedToFormatField](&FailedToFormatField{})(nil),
	}

	// CNAB_ErrFieldExceedsSpecifiedLength is an error that occurs when a field exceeds the specified length.
	CNAB_ErrFieldExceedsSpecifiedLength = &ErrInstance[*FieldExceedsSpecifiedLength]{
		Encapsulator: &FieldExceedsSpecifiedLength{},
		Creator:      iError.NewCreator[*FieldExceedsSpecifiedLength](&FieldExceedsSpecifiedLength{}),
		Err:          iError.NewCreator[*FieldExceedsSpecifiedLength](&FieldExceedsSpecifiedLength{})(nil),
	}

	// CNAB_ErrFieldIsEmpty is an error that occurs when a field is empty.
	CNAB_ErrFieldIsEmpty = &ErrInstance[*FieldIsEmpty]{
		Encapsulator: &FieldIsEmpty{},
		Creator:      iError.NewCreator[*FieldIsEmpty](&FieldIsEmpty{}),
		Err:          iError.NewCreator[*FieldIsEmpty](&FieldIsEmpty{})(nil),
	}

	// CNAB_ErrUnsupportedFieldType is an error that occurs when a field type is unsupported.
	CNAB_ErrUnsupportedFieldType = &ErrInstance[*UnsupportedFieldType]{
		Encapsulator: &UnsupportedFieldType{},
		Creator:      iError.NewCreator[*UnsupportedFieldType](&UnsupportedFieldType{}),
		Err:          iError.NewCreator[*UnsupportedFieldType](&UnsupportedFieldType{})(nil),
	}

	// CNAB_ErrFieldValueIsNotAnDate is an error that occurs when a field value is not a date.
	CNAB_ErrFieldValueIsNotAnDate = &ErrInstance[*FieldValueIsNotAnDate]{
		Encapsulator: &FieldValueIsNotAnDate{},
		Creator:      iError.NewCreator[*FieldValueIsNotAnDate](&FieldValueIsNotAnDate{}),
		Err:          iError.NewCreator[*FieldValueIsNotAnDate](&FieldValueIsNotAnDate{})(nil),
	}

	// CNAB_ErrFieldValueIsNotAnString is an error that occurs when a field value is not a string.
	CNAB_ErrFieldValueIsNotAnString = &ErrInstance[*FieldValueIsNotAnString]{
		Encapsulator: &FieldValueIsNotAnString{},
		Creator:      iError.NewCreator[*FieldValueIsNotAnString](&FieldValueIsNotAnString{}),
		Err:          iError.NewCreator[*FieldValueIsNotAnString](&FieldValueIsNotAnString{})(nil),
	}

	// CNAB_ErrFieldValueIsNotAnInt is an error that occurs when a field value is not an int.
	CNAB_ErrFieldValueIsNotAnInt = &ErrInstance[*FieldValueIsNotAnInt]{
		Encapsulator: &FieldValueIsNotAnInt{},
		Creator:      iError.NewCreator[*FieldValueIsNotAnInt](&FieldValueIsNotAnInt{}),
		Err:          iError.NewCreator[*FieldValueIsNotAnInt](&FieldValueIsNotAnInt{})(nil),
	}

	// CNAB_ErrFieldValueIsNotAnFloat is an error that occurs when a field value is not a float.
	CNAB_ErrFieldValueIsNotAnFloat = &ErrInstance[*FieldValueIsNotAnFloat]{
		Encapsulator: &FieldValueIsNotAnFloat{},
		Creator:      iError.NewCreator[*FieldValueIsNotAnFloat](&FieldValueIsNotAnFloat{}),
		Err:          iError.NewCreator[*FieldValueIsNotAnFloat](&FieldValueIsNotAnFloat{})(nil),
	}

	// CNAB_ErrFieldExceedsRecordLength is an error that occurs when a field exceeds the record length.
	CNAB_ErrFieldExceedsRecordLength = &ErrInstance[*FieldExceedsRecordLength]{
		Encapsulator: &FieldExceedsRecordLength{},
		Creator:      iError.NewCreator[*FieldExceedsRecordLength](&FieldExceedsRecordLength{}),
		Err:          iError.NewCreator[*FieldExceedsRecordLength](&FieldExceedsRecordLength{})(nil),
	}

	// CNAB_ErrFailedToParseField is an error that occurs when a field cannot be parsed.
	CNAB_ErrFailedToParseField = &ErrInstance[*FailedToParseField]{
		Encapsulator: &FailedToParseField{},
		Creator:      iError.NewCreator[*FailedToParseField](&FailedToParseField{}),
		Err:          iError.NewCreator[*FailedToParseField](&FailedToParseField{})(nil),
	}

	// CNAB_ErrCannotConvertToInt is an error that occurs when a value cannot be converted to an int.
	CNAB_ErrCannotConvertToInt = &ErrInstance[*CannotConvertToInt]{
		Encapsulator: &CannotConvertToInt{},
		Creator:      iError.NewCreator[*CannotConvertToInt](&CannotConvertToInt{}),
		Err:          iError.NewCreator[*CannotConvertToInt](&CannotConvertToInt{})(nil),
	}

	// CNAB_ErrCannotConvertToFloat is an error that occurs when a value cannot be converted to a float.
	CNAB_ErrCannotConvertToFloat = &ErrInstance[*CannotConvertToFloat]{
		Encapsulator: &CannotConvertToFloat{},
		Creator:      iError.NewCreator[*CannotConvertToFloat](&CannotConvertToFloat{}),
		Err:          iError.NewCreator[*CannotConvertToFloat](&CannotConvertToFloat{})(nil),
	}

	// CNAB_ErrInvalidDecimalValue is an error that occurs when a decimal value is invalid.
	CNAB_ErrInvalidDecimalValue = &ErrInstance[*InvalidDecimalValue]{
		Encapsulator: &InvalidDecimalValue{},
		Creator:      iError.NewCreator[*InvalidDecimalValue](&InvalidDecimalValue{}),
		Err:          iError.NewCreator[*InvalidDecimalValue](&InvalidDecimalValue{})(nil),
	}

	// CNAB_ErrMissingDateFormat is an error that occurs when a date format is missing.
	CNAB_ErrMissingDateFormat = &ErrInstance[*MissingDateFormat]{
		Encapsulator: &MissingDateFormat{},
		Creator:      iError.NewCreator[*MissingDateFormat](&MissingDateFormat{}),
		Err:          iError.NewCreator[*MissingDateFormat](&MissingDateFormat{})(nil),
	}

	// CNAB_ErrInvalidDateLength is an error that occurs when a date length is invalid.
	CNAB_ErrInvalidDateLength = &ErrInstance[*InvalidDateLength]{
		Encapsulator: &InvalidDateLength{},
		Creator:      iError.NewCreator[*InvalidDateLength](&InvalidDateLength{}),
		Err:          iError.NewCreator[*InvalidDateLength](&InvalidDateLength{})(nil),
	}
)
