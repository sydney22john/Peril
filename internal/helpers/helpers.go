package helpers

import (
	"os"
	"os/signal"
)

func BlockUntilSignal(sig os.Signal) {
	// wait for signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, sig)
	<-signalChan
}
