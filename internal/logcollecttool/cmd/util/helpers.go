package util

import (
	"bytes"
	"fmt"
	"logconvert/errors"
	"net/url"
	"os"
	"strings"

	"logconvert/log"
)

const (
	// DefaultErrorExitCode defines the default exit code.
	DefaultErrorExitCode = 1
)

type debugError interface {
	DebugError() (msg string, args []interface{})
}

var fatalErrHandler = fatal

// fatal prints the message (if provided) and then exits.
func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

// ErrExit may be passed to CheckError to instruct it to output nothing but exit with
// status code 1.
var ErrExit = fmt.Errorf("exit")

func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

func checkErr(err error, handleErr func(string, int)) {
	// unwrap aggregates of 1
	if agg, ok := err.(errors.Aggregate); ok && len(agg.Errors()) == 1 {
		err = agg.Errors()[0]
	}

	if err == nil {
		return
	}

	switch {
	case err == ErrExit:
		handleErr("", DefaultErrorExitCode)
	default:
		switch err := err.(type) {
		case errors.Aggregate:
			handleErr(MultipleErrors(``, err.Errors()), DefaultErrorExitCode)
		default: // for any other error type
			msg, ok := StandardErrorMessage(err)
			if !ok {
				msg = err.Error()
				if !strings.HasPrefix(msg, "error: ") {
					msg = fmt.Sprintf("error: %s", msg)
				}
			}
			handleErr(msg, DefaultErrorExitCode)
		}
	}
}

func StandardErrorMessage(err error) (string, bool) {
	if debugErr, ok := err.(debugError); ok {
		log.Infof(debugErr.DebugError())
	}
	if t, ok := err.(*url.Error); ok {
		log.Infof("Connection error: %s %s: %v", t.Op, t.URL, t.Err)
		if strings.Contains(t.Err.Error(), "connection refused") {
			host := t.URL
			if server, err := url.Parse(t.URL); err == nil {
				host = server.Host
			}
			return fmt.Sprintf(
				"The connection to the server %s was refused - did you specify the right host or port?",
				host,
			), true
		}

		return fmt.Sprintf("Unable to connect to the server: %v", t.Err), true
	}
	return "", false
}

// MultipleErrors returns a newline delimited string containing
// the prefix and referenced errors in standard form.
func MultipleErrors(prefix string, errs []error) string {
	buf := &bytes.Buffer{}
	for _, err := range errs {
		fmt.Fprintf(buf, "%s%v\n", prefix, messageForError(err))
	}
	return buf.String()
}

// messageForError returns the string representing the error.
func messageForError(err error) string {
	msg, ok := StandardErrorMessage(err)
	if !ok {
		msg = err.Error()
	}
	return msg
}

// ValidateOptions defines the validate options.
type ValidateOptions struct {
	EnableValidation bool
}
