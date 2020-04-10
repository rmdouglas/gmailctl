package errors

import (
	"errors"
	"fmt"
	"io"
)

// Aliases to the standard errors package.
var (
	New = errors.New
	Is  = errors.Is
	As  = errors.As
)

// New is an alias for errors.New.

// WithCause annotates a symptom error with a cause.
//
// Both errors can be discovered by the Is and As methods.
func WithCause(symptom, cause error) error {
	return annotated{
		cause:   cause,
		symptom: symptom,
	}
}

func WithDetails(err error, details string) error {
	if err == nil {
		return nil
	}
	return detailed{err, details}
}

func Details(err error) string {
	var dErr detailed
	if errors.As(err, &dErr) {
		return dErr.details
	}
	return ""
}

type detailed struct {
	error
	details string
}

func (e detailed) Unwrap() error {
	return e.error
}

func (w detailed) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.error)
			io.WriteString(s, w.details)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

type annotated struct {
	cause   error
	symptom error
}

func (e annotated) Error() string {
	return fmt.Sprintf("%s: %s", e.cause, e.symptom)
}

func (e annotated) Unwrap() error {
	return e.cause
}

func (e annotated) Is(target error) bool {
	return errors.Is(e.symptom, target) || errors.Is(e.cause, target)
}

func (e annotated) As(target interface{}) bool {
	if errors.As(e.symptom, target) {
		return true
	}
	return errors.As(e.cause, target)
}
