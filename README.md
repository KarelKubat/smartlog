# Smartlog

Smartlog is a package for Go to make setting up logging easier. Log statements can be processed locally (to `stdout` or a file) or sent remotely to a server over TCP or UDP for further handling.

Smartlog contains all support code to embed such logging into your programs:

- To embed client code into your programs that need to emit log messages
- To embed server code into a centralized server for further processing.

## Concepts
---

### Emitting messages from your Go program

A program that wishes to provide some logging information uses a Smartlog client to emit messages. Smartlog supports several message types:

- Debug messages, which are emitted when a debug level is exceeded. You can sprinkle calls to `client.Debug(lev, msg)` with different levels in your program and then set an appropriate threshold to either have these emitted or suppressed.
- Informational messages: `client.Info(msg)`
- Warnings: `client.Warn(msg)`
- Fatal errors: `client.Fatal(msg)` which also cause the program to exit.
- There are corresponding `~f()` versions that accept a format string and arguments, a-la `fmt.Printf()` - e.g., `client.Warnf(format, ...args)`.

There are tons of discussions on what logging should be aimed at, what it should do, and especially what it should not do. Smartlog isn't as pure as the suggestions by [Dave Cheney](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) but instead chooses the following approach:

- Debug messages can be used during development and should be aimed at programmers. You can leave them in the code; in production they can be turned into no-ops by choosing an appropriate level. Or, if needed, you can turn up the level and see what's going on.
- Informational messages are aimed at users in order to provide relevant (business) data, like "your bank balance looks great today".
- Warnings are just informational messages that should stand out, like "your bank balance is dangerously low". They don't fix anything; the dangerous situation still needs to be handled by your program.
- Fatals should not be used, except in the simplest of programs where it's ok to `exit(1)` and to abandon all running threads, pending file writes, etc.. Programs that need cleanups should just issue a warning, and let the appropriate error bubble up to `main()` for handling.

Smartlog servers have a queue for incoming messages. When this queue fills up (i.e., messages are received faster than they are handled) then debug messages are discarded first. If the queue still fills up, informational messages are discarded. Received arnings and fatals are never discarded.

### Client types

Client types define how a message should be handled. Smartlog supports the following types:

- File-based clients dump messages into a file. New clients that point to the same file append to the file instead of overwriting it (this is also the case when you re-run your program and point to the same file as the last time). The file may disappear while your program is running; in that case, smartlog will simply re-open it. This is practical for logfile rotation: an external script may e.g. move the file and zip it, and smartlog will simply create a new one.
- A special case is the filename `stdout`, which instructs smartlog to send messges to the console.
- Network-based clients send messages to a remote server. Smartlog supports UDP and TCP:
  - UDP is faster, but the network transmission is not guaranteed.
  - TCP is slower, but guaranteed.

## Tweaks
---

### Timestamps

Any client-side invocation like `client.Info("hello world")` leads to a message which has the timestamp. Two settings can be controlled:

- The timestamp format: the default is `"2006-01-02 15:04:05 MST"` (see e.g. the [Go time package](https://pkg.go.dev/time) or [Geeks for geeks](https://www.geeksforgeeks.org/time-formatting-in-golang/)
- Whether the time is displayed relative to localtime or to UTC: the default is `false`: the localtime is shown, not the UTC time.

To change the defaults, simply modify the global variables in `smartlog/msg`:

```go
import (
    "time"
    "smartlog/msg"
)    
// ...
msg.DefaultTimeFormat = time.RFC3339 // format: "2006-01-02T15:04:05Z07:00"
msg.UTCTime = true                   // show the UTC time, not the localtime
```