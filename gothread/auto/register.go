//Package auto listens for SIGQUIT signals and prints a threaddunp to Stderr.
//This is done via an init function and expected to be registered by an _ import in main
package auto

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/devmop/gothread"
)

func init() {
	ch := make(chan os.Signal, 10)
	go listen(ch)
	signal.Notify(ch, syscall.SIGQUIT)
}

func listen(ch chan os.Signal) {
	for _ = range ch {
		gothread.WriteStackTrace(os.Stderr)
	}
}
