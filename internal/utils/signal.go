package utils

import (
	"os"
	"os/signal"
)

func BlockUntilSignal() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
}
