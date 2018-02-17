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
        for i := 0; running ; i++ {
            k, ks, err := l.Pop()
            fmt.Printf("%+v: %+v (err: %+v\n", k, ks, err)
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
