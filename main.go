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

    go func() {
        for i := 0; ; i++ {
            fmt.Println(i, "hey...")
            time.Sleep(time.Millisecond * 500)
        }
    } ()

    l.Pop()
    l.Wait()
    return
}
