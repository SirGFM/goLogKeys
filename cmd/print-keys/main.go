package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/SirGFM/goLogKeys/logger"
)

func main() {
	l, err := logger.GetLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to initialized the logger: %+v", err))
	}

	defer l.Clean()
	err = l.Setup()
	if err != nil {
		panic(fmt.Sprintf("failed to setup the logger: %+v", err))
	}

	err = l.Start()
	if err != nil {
		panic(fmt.Sprintf("failed to start the logger: %+v", err))
	}

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	intHndlr := make(chan os.Signal, 1)
	signal.Notify(intHndlr, os.Interrupt)

	keys := make([]logger.Key, 100)
	states := make([]logger.KeyState, 100)

	for {
		select {
		case <-t.C:
			k, s, err := l.PopMulti(keys, states)
			if err != nil {
				panic(fmt.Sprintf("failed to read keys: %+v", err))
			}

			for i := range k {
				fmt.Printf("%s -> %s\n", k[i], s[i])
			}
		case <-intHndlr:
			return
		}
	}
}
