package analyzer

import (
	"errors"
	"fmt"
)

type ErrorCode string

func (c ErrorCode) String() string {
	return string(c)
}

const (
	CodeInvalidReference ErrorCode = "INVALID_REFERENCE"
	CodeImageNotFound    ErrorCode = "IMAGE_NOT_FOUND"
	CodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	CodeTimeout          ErrorCode = "TIMEOUT"
	CodeNoLayers         ErrorCode = "NO_LAYERS"
	CodeFetchFailed      ErrorCode = "FETCH_FAILED"
	CodeBuildFailed      ErrorCode = "BUILD_FAILED"
	CodeDigestFailed     ErrorCode = "DIGEST_FAILED"
	CodeMediaTypeFailed  ErrorCode = "MEDIA_TYPE_FAILED"
	CodeSizeFailed       ErrorCode = "SIZE_FAILED"
	CodeLayerExtract     ErrorCode = "LAYER_EXTRACT_FAILED"
	CodeUnknown          ErrorCode = "UNKNOWN"
)

type AnalyzerError struct {
	code    ErrorCode
	message string
	err     error
	op      string
	ref     string
}

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

func (e *AnalyzerError) Unwrap() error {
	return e.err
}

// Allow errors.Is(err, target)
func (e *AnalyzerError) Is(target error) bool {
	t, ok := target.(*AnalyzerError)
	if !ok {
		return false
	}
	return e.code == t.code
}

func (e *AnalyzerError) Code() ErrorCode {
	return e.code
}

func (e *AnalyzerError) Operation() string {
	return e.op
}

func (e *AnalyzerError) Reference() string {
	return e.ref
}

func (e *AnalyzerError) Temporary() bool {
	return e.code == CodeTimeout || e.code == CodeFetchFailed
}

func (e *AnalyzerError) Timeout() bool {
	return e.code == CodeTimeout
}

func NewError(code ErrorCode, op, ref, message string, err error) *AnalyzerError {
	return &AnalyzerError{
		code:    code,
		message: message,
		err:     err,
		op:      op,
		ref:     ref,
	}
}

func Wrap(code ErrorCode, op, ref string, err error) *AnalyzerError {
	return NewError(code, op, ref, code.String(), err)
}

/*
	Sentinel base errors (optional for compatibility)
*/

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

func IsCode(err error, code ErrorCode) bool {
	var ae *AnalyzerError
	if errors.As(err, &ae) {
		return ae.code == code
	}
	return false
}

func AsAnalyzerError(err error) (*AnalyzerError, bool) {
	var ae *AnalyzerError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
