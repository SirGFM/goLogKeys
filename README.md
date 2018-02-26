# Go Log Keys

Simple key logger written in C and Go.

Its intended use it to register some keys (from a JSON file or something), and
register key events (press/release) through a WebSocket to my
[stream overlay](https://github.com/SirGFM/Roa-Stream-Skin), so it can display
which keyboard keys (mapped to a controller or something) are currently pressed.


## Logging on Windows

Windows allows setting a callback for global keyboard events. However, it must
have been loaded from a DLL, instead of from the current process...

`./winDll/logger.c` implements the callback, `keyCallback`, and a function
to retrieve all keypresses, `pop`.

Note that (from this
[StackOverflow answer](https://stackoverflow.com/questions/11176408/setwindowshookex-with-wh-keyboard-doesnt-work-for-me-what-do-i-wrong))
that `sleep` etc blocks the callback... So, that DLL also has a function,
`wait`, that creates a dummy window, which may be used to keep the logger alive.

With [mxe](http://mxe.cc/) installed, you may compile the DLL and the test with:

```sh
export CC=x86_64-w64-mingw32.shared-gcc
${CC} -shared -Wl,-soname,logger.dll -Wl,-export-all-symbols -o logger64.dll winDll/logger.c
${CC} -o test64.exe winDll/test.c
```

Note that I use the 64-bit MinGW compiler... I could only get this to work by
using a 64-bit DLL along with a 64-bit test application (in my 64-bit system).
However, I believe that injecting a 32-bit DLL into a 32-bit application in a
64-bit system should be possible...

Building the Go part should be quite straight forward... Try to build it and
`go get` anything that fails:

```sh
go get -d github.com/pkg/errors
go get -d golang.org/x/sys/windows
go get -d github.com/SirGFM/GoWebSocketProxy/websocket
GOOS=windows GOARCH=amd64 go build
```


## Logging on Linux

TODO D:

Should be easy, since `key-mon` does exactly this...


## Resources

Non exhaustive list of resources used to write all this:

* [MSDN: LowLevelKeyboardProc callback function](https://msdn.microsoft.com/en-US/library/windows/desktop/ms644985(v=vs.85).aspx)
* [MSDN: KBDLLHOOKSTRUCT structure](https://msdn.microsoft.com/en-US/library/windows/desktop/ms644967(v=vs.85).aspx)
* [MSDN: SetWindowsHookEx function](https://msdn.microsoft.com/en-US/library/windows/desktop/ms644990(v=vs.85).aspx)
* [MSDN: Installing and Releasing Hook Procedures](https://msdn.microsoft.com/en-US/library/windows/desktop/ms644960(v=vs.85).aspx#installing_releasing)
* [SO: SetWindowsHookEx with WH_KEYBOARD doesn't work for me, what do I wrong?](https://stackoverflow.com/questions/11176408/setwindowshookex-with-wh-keyboard-doesnt-work-for-me-what-do-i-wrong)
* [Go: Calling a Windows DLL](https://github.com/golang/go/wiki/WindowsDLLs)
* [golang.org/x/sys/windows](https://godoc.org/golang.org/x/sys/windows)

