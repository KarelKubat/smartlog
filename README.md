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
- Errors: `client.Error(msg)` which also cause the program to exit.

There are corresponding `-f()` versions that accept a format string and arguments, a-la `fmt.Printf()` - e.g., `client.Warnf(format, ...args)`.

Smartlog clients have a message queue for emitting. When this queue fills up (i.e., messages are generated faster than they are handled) then debug messages are discarded first. If the queue still fills up, informational messages are discarded. Warnings and errors are never discarded.

### Client types

Client types define how a message should be handled. Smartlog supports the following types:

- File-based clients dump messages into a file. New clients that point to the same file append to the file instead of overwriting it (this is also the case when you re-run your program and point to the same file as the last time). The file may disappear while your program is running; in that case, smartlog will simply re-open it. This is practical for logfile rotation: an external script may e.g. move the file and zip it, and smartlog will simply create a new one.
- A special case is the filename `stdout`, which instructs smartlog to send messges to the console.
- Network-based clients send messages to a remote server. Smartlog supports UDP and TCP:
  - UDP is faster, but the network transmission is not guaranteed.
  - TCP is slower, but guaranteed.

## Smartlog Clients
---

### Timestamps

Any client-side invocation like `client.Info("hello world")` leads to a message which has the timestamp. Two settings can be controlled:

- The timestamp format: the default is `"2006-01-02 15:04:05 MST"` (see e.g. https://www.geeksforgeeks.org/time-formatting-in-golang/)
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