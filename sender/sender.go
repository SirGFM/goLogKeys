package sender

import (
    "github.com/SirGFM/GoWebSocketProxy/websocket"
    "github.com/SirGFM/goLogKeys/logger"
    "time"
)

// How many updates should be sent per second
const updateRate = 30
// Delay between each update
const updateDelay = time.Second / updateRate

type server struct {
    running bool
    conn websocket.ClientConnection
    timer *time.Ticker

    // The main logger, shared between every server instance (and copied from
    // the main instance).
    loggerCtx logger.Logger
    // Filter keys accepted by the sender and map them into their a button bitmask.
    filter map[logger.Key]uint32
}

// Retrieve a new sender template, with the shared logger and the map to filter
// accepted keys.
func GetTemplate(loggerCtx logger.Logger, acceptedKeys []logger.Key) (
    websocket.Server) {

    if len(acceptedKeys) == 0 {
        panic("goLogKeys.sender: Missing keys to be filtered")
    } else if len(acceptedKeys) > 32 {
        panic("goLogKeys.sender: Too many keys to be filtered")
    }

    filter := make(map[logger.Key]uint32)
    i := uint32(1)
    for _, key := range acceptedKeys {
        filter[key] = i
        i <<= 1
    }

    return &server{
        loggerCtx: loggerCtx,
        filter: filter,
    }
}

// Start a new timer for the connection, and the goroutine that shall send the
// inputs to the client.
func (s *server) Clone(conn websocket.ClientConnection) (websocket.Server, error) {
    t := time.NewTicker(updateDelay)

    newServer := server{
        running: true,
        conn:    conn,
        timer:   t,
        loggerCtx:   s.loggerCtx,
        filter:   s.filter,
    }

    // Start sending data
    go newServer.run()

    return &newServer, nil
}

// Release resources associated with the timer and stops the goroutine.
func (srv *server) Cleanup() {
    srv.running = false
    time.Sleep(updateDelay * 2)
    srv.timer.Stop()
}

// The client should never send any data, so simply ignore it...
func (*server) Do(msg []byte, offset int) (err error) {
    return
}

// Lookup table used to slightly speed up converting an integer value into a
// string.
const _itoaLookup = "0123456789ABCDEF"

// Converts a bitmask (store as a uint32) into a byte slice. The value is
// written in big-endian format, as a hexadecimal number.
func itoa(buf []byte, mask uint32) (out []byte) {
    out = buf[:0]

    for ; mask != 0; mask >>= 4 {
        out = append(out, _itoaLookup[mask & 0xF])
    }

    return
}

// Keep sending data to the client.
func (srv *server) run() {
    var keys []logger.Key
    var states []logger.KeyState
    var tmpBuf []byte

    for {
        var err error

        state := uint32(0)

        // Wait for a new update
        _ = <-srv.timer.C
        if !srv.running {
            break
        }

        // Read every key press since the last update
        for {
            keys, states, err = srv.loggerCtx.PopMulti(keys, states)
            if err != nil {
                // Failed to read the keys... Who know why? There isn't much
                // else that could be done after this. D:
                panic(err.Error())
            } else if len(keys) == 0 {
                break
            }

            // Filter only the desired keys
            for i, s := range states {
                // If the key isn't accepted by the logger (i.e., if it isn't in
                // the map), the map returns 0.
                bit := srv.filter[keys[i]]

                if s == logger.Pressed {
                    state |= bit
                } else {
                    state &^= bit
                }
            }
        }

        // Send the packed keys
        tmpBuf = itoa(tmpBuf, state)
        _, err = srv.conn.Write(tmpBuf)
        if err != nil {
            // Failed to send more that... This isn't as bad as failling to read
            // keys, but this connection to the client may not be realiable
            // anymore. Close this connection and wait for a new one.
            err = srv.conn.Close(websocket.UnexpectedError, []byte(err.Error()))
            if err != nil {
                // Because of the dumb way that I wrote the WebSocket server, if
                // this fails there wouldn't be much else to be done (again!)...
                // So, yeah...
                panic(err.Error())
            }
            return
        }
    }
}
