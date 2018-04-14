package detailerror

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
)

// This is implement of error interface.
type DetailError struct {
	ErrNo  int
	ErrMsg string
	ErrPos string
}

func (e *DetailError) Error() string {
	return e.String()
}

func (e *DetailError) String() string {
	if e == nil {
		return fmt.Sprint("[errno:-1 errmsg:nil]")
	}
	return fmt.Sprintf("[errno:%d errmsg:%s errpos:%s]", e.ErrNo, e.ErrMsg, e.ErrPos)
}

func newDetailError(errno int, errmsg string) *DetailError {
	var errpos string
	_, file, line, ok := runtime.Caller(2)
	if ok {
		errpos = path.Base(file) + ":" + strconv.Itoa(line)
	}
	return &DetailError{errno, errmsg, errpos}
}

func New(errno int, errmsg string) *DetailError {
	return newDetailError(errno, errmsg)
}

func Errorf(errno int, format string, a ...interface{}) *DetailError {
	return newDetailError(errno, fmt.Sprintf(format, a...))
}
