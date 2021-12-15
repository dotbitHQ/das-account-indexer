package toolib

import (
	"os"
	"os/signal"
	"syscall"
)

func ExitMonitoring(handle func(sig os.Signal)) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)                   // signal int, kill session windows
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT) // ctrl+C
	signal.Notify(ch, os.Kill, syscall.SIGKILL)       // kill -9
	signal.Notify(ch, syscall.SIGTERM)                // kill -15
	signal.Notify(ch, syscall.SIGQUIT)
	go func() {
		for {
			select {
			case s := <-ch:
				handle(s)
			}
		}
	}()
}
