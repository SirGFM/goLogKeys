package main

import (
    "github.com/SirGFM/goLogKeys/logger"

    "fmt"
    "time"
)

func main() {
    l, err := logger.GetLogger()
    if err != nil {
        panic(err.Error())
    }
    defer l.Clean()

    if err := l.Setup(); err != nil {
        panic(err.Error())
    }

    running := true
    go func() {
        var keys []logger.Key
        var states []logger.KeyState

        for i := 0; running ; i++ {
            var err error

            keys, states, err = l.PopMulti(keys, states)
            if err != nil {
                fmt.Printf("Error: %+v\n", err)
            } else {
                for i := range keys {
                    k := keys[i]
                    ks := states[i]
                    fmt.Printf("%+v: %+v\n", k, ks)
                }
            }

            time.Sleep(time.Millisecond * 500)
        }
    } ()

    if err := l.Start(); err != nil {
        panic(err.Error())
    }

    l.Wait()
    running = false
    return
}
