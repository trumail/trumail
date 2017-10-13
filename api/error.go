package api

import (
	"encoding/xml"
	"fmt"

	"github.com/Sirupsen/logrus"
)

// Error is a generic API error
type Error struct {
	XMLName    xml.Name `json:"-" xml:"error"`
	Message    string   `json:"message,omitempty" xml:"message,omitempty"`
	StatusCode int      `json:"status,omitempty" xml:"status,omitempty"`
	Err        string   `json:"err,omitempty" xml:"err,omitempty"`
}

// NewError is a generic http error type used for all error responses
func NewError(message string, statusCode int, err error) *Error {
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &Error{
		Message:    message,
		StatusCode: statusCode,
		Err:        errString,
	}
}

// Error returns a string representation of the error and
// helps to satisfy the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("Error %d: '%s'", e.StatusCode, e.Message)
}

// Log will log the Error before returning it
func (e *Error) Log(log *logrus.Entry) *Error {
	log.WithFields(logrus.Fields{"error": e.Err, "status": e.StatusCode}).Error(e.Message)
	return e
}
