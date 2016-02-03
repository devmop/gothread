package gothread

import (
	"io"
	"runtime/pprof"
)

//WriteStackTrace atempts to write the current state of all goroutines to the provided writer
func WriteStackTrace(w io.Writer) error {
	return pprof.Lookup("goroutine").WriteTo(w, 2)
}
