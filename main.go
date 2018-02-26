package main

import (
    "encoding/json"
    "fmt"
    "github.com/SirGFM/GoWebSocketProxy/websocket"
    "github.com/SirGFM/goLogKeys/logger"
    "github.com/SirGFM/goLogKeys/sender"
    "io"
    "os"
)

type Args struct {
    // Only URL (if any) from which the server may be accessed.
    Url string
    // URI where the server is accessed.
    Uri string
    // Listening port of the server.
    Port int
    // List of keys that the server filters and pass on to the client. The keys
    // are assigned bits sequentially.
    Keys []string
    // NOTE: This field is filled from the information found on the JSON file.
    acceptedKeys []logger.Key
}

// Parse the input file as a JSON. The file must not have anything other then a
// single 'Args' object.
func parseArgs(filePath string) *Args {
    var args Args

    f, err := os.Open(filePath)
    if err != nil {
        panic(fmt.Sprintf("Failed to open the input file: %+v", err))
    }
    defer f.Close()

    dec := json.NewDecoder(f)

    err = dec.Decode(&args)
    if err != nil {
        panic(fmt.Sprintf("Failed to parse the input file: %+v", err))
    }

    _, err = dec.Token()
    if err != io.EOF {
        panic("Found trailing data at the end of the JSON file")
    }

    for _, k := range args.Keys {
        args.acceptedKeys = append(args.acceptedKeys, logger.GetKey(k))
    }

    return &args
}

func main() {
    // Parse the configurations
    if len(os.Args) != 2 {
        panic(fmt.Sprintf("Invalid number of arguments! Usage: %s filter.json", os.Args[0]))
    }
    args := parseArgs(os.Args[1])

    // Prepare the key logger
    keyLogger, err := logger.GetLogger()
    if err != nil {
        panic(fmt.Sprintf("Failed to setup the logger: %+v", err))
    }
    defer keyLogger.Clean()

    // Create the WebSocket server (but do not start handling it just yet)
    wsCtx := websocket.NewContext(args.Url, args.Uri, args.Port, 1)
    defer wsCtx.Close()
    err = wsCtx.Setup(sender.GetTemplate(keyLogger, args.acceptedKeys))
    if err != nil {
        panic(fmt.Sprintf("Failed to setup the WebSocket server: %+v", err))
    }

    // Start logging keys
    err = keyLogger.Setup()
    if err != nil {
        panic(fmt.Sprintf("Failed to start logging keys: %+v", err))
    }
    err = keyLogger.Start()
    if err != nil {
        panic(fmt.Sprintf("Failed to start logging keys: %+v", err))
    }

    // Start the WebSocket server
    go func (ctx *websocket.Context) {
        cerr := make(chan error)
        go ctx.Run(cerr)

        for {
            err := <-cerr
            if err == nil {
                break
            }

            fmt.Printf("Got error from server: %+v\n", err)
        }
    } (wsCtx)

    // Wait until the application is closed
    keyLogger.Wait()
}
