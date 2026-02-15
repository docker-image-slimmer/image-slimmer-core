package analyzer

import (
	"errors"
	"fmt"
)

// ErrorCode represents a machine-readable classification of an analyzer error.
type ErrorCode string

// String returns the string representation of the error code.
func (c ErrorCode) String() string {
	return string(c)
}

const (
	// Reference parsing errors
	CodeInvalidReference ErrorCode = "INVALID_REFERENCE"

	// Registry and fetch errors
	CodeImageNotFound ErrorCode = "IMAGE_NOT_FOUND"
	CodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	CodeTimeout       ErrorCode = "TIMEOUT"
	CodeFetchFailed   ErrorCode = "FETCH_FAILED"

	// Image structure errors
	CodeNoLayers        ErrorCode = "NO_LAYERS"
	CodeBuildFailed     ErrorCode = "BUILD_FAILED"
	CodeDigestFailed    ErrorCode = "DIGEST_FAILED"
	CodeMediaTypeFailed ErrorCode = "MEDIA_TYPE_FAILED"
	CodeSizeFailed      ErrorCode = "SIZE_FAILED"
	CodeLayerExtract    ErrorCode = "LAYER_EXTRACT_FAILED"

	CodeValidationFailed ErrorCode = "VALIDATION_FAILED"

	// Fallback classification
	CodeUnknown ErrorCode = "UNKNOWN"
)

// AnalyzerError is the unified structured error type used across the analyzer module.
// It encapsulates classification, context and the underlying cause.
type AnalyzerError struct {
	code    ErrorCode // machine-readable classification
	message string    // human-readable description
	err     error     // wrapped underlying error
	op      string    // logical operation name (e.g., "fetch", "build")
	ref     string    // image reference involved in the error
}

// Error implements the error interface.
func (e *AnalyzerError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%s] %s (op=%s ref=%s): %v",
			e.code,
			e.message,
			e.op,
			e.ref,
			e.err,
		)
	}

	return fmt.Sprintf("[%s] %s (op=%s ref=%s)",
		e.code,
		e.message,
		e.op,
		e.ref,
	)
}

// Unwrap allows errors.Unwrap / errors.Is / errors.As to access the underlying error.
func (e *AnalyzerError) Unwrap() error {
	return e.err
}

// Is enables errors.Is comparison based on error code equality.
func (e *AnalyzerError) Is(target error) bool {
	t, ok := target.(*AnalyzerError)
	if !ok {
		return false
	}
	return e.code == t.code
}

// Code returns the structured error classification.
func (e *AnalyzerError) Code() ErrorCode {
	return e.code
}

// Operation returns the logical operation that produced the error.
func (e *AnalyzerError) Operation() string {
	return e.op
}

// Reference returns the image reference associated with the error.
func (e *AnalyzerError) Reference() string {
	return e.ref
}

// Temporary reports whether the error is potentially retryable.
func (e *AnalyzerError) Temporary() bool {
	return e.code == CodeTimeout || e.code == CodeFetchFailed
}

// Timeout reports whether the error represents a timeout condition.
func (e *AnalyzerError) Timeout() bool {
	return e.code == CodeTimeout
}

// NewError creates a new structured AnalyzerError.
func NewError(code ErrorCode, op, ref, message string, err error) *AnalyzerError {
	return &AnalyzerError{
		code:    code,
		message: message,
		err:     err,
		op:      op,
		ref:     ref,
	}
}

// Wrap creates a new AnalyzerError using the error code string as default message.
func Wrap(code ErrorCode, op, ref string, err error) *AnalyzerError {
	return NewError(code, op, ref, code.String(), err)
}

/*
	Sentinel base errors (optional for compatibility)
*/

// These sentinel errors allow usage with errors.Is without requiring full struct matching.
var (
	ErrInvalidReference = &AnalyzerError{code: CodeInvalidReference}
	ErrImageNotFound    = &AnalyzerError{code: CodeImageNotFound}
	ErrUnauthorized     = &AnalyzerError{code: CodeUnauthorized}
	ErrTimeout          = &AnalyzerError{code: CodeTimeout}
	ErrNoLayers         = &AnalyzerError{code: CodeNoLayers}
	ErrFetchFailed      = &AnalyzerError{code: CodeFetchFailed}
	ErrBuildFailed      = &AnalyzerError{code: CodeBuildFailed}
)

/*
	Helpers
*/

// IsCode checks whether an error matches a specific ErrorCode.
func IsCode(err error, code ErrorCode) bool {
	var ae *AnalyzerError
	if errors.As(err, &ae) {
		return ae.code == code
	}
	return false
}

// AsAnalyzerError attempts to extract an AnalyzerError from a generic error.
func AsAnalyzerError(err error) (*AnalyzerError, bool) {
	var ae *AnalyzerError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
