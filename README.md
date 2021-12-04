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
- HTTP clients start an HTTP server where messages can be viewed.
- Forwarding clients send messages to a remote server. Smartlog supports UDP and TCP:
  - UDP is faster, but the network transmission is not guaranteed.
  - TCP is slower, but guaranteed.

All client types except the forwarding clients can be used stand-alone, i.e., just as a part of your program. Forwarding clients require a Smartlog server.

### Smartlog servers need smartlog clients too

A Smartlog server (which receives messages over TCP or UDP) is in itself useless. It needs clients to do something with incoming messages. The clients that a server uses are are identical to any client that you'd use in your own program: messages arriving at the server may be sent to a file, to `stdout`, kept for viewing in an HTTP client, or forwarded to next hops (and the story repeats at the Smartlog servers that accept those messages).

Here is an example that uses ready-to-run programs in the package:

1. In one terminal run:

   ```sh
   # Terminal #1
   # The first positional argument is the server, others are clients.
   # Accept messages on TCP, port 2022. Save them to /tmp/out.txt and make them viewable on http://localhost:8080.
   # tcp://:2022 means any IP on this machine. The filename /tmp/out.txt leads to three slashes in
   # file:///tmp/out.txt; file:// already needs 2.
   go run main/server/smartlog-server.go tcp://:2022 file:///tmp/out.txt http://localhost:8080
   ```

1. In another terminal run:

   ```sh
   # Terminal #2
   # Accept messages on UDP, port 2021. Fan these out to `stdout` and forward them to the TCP server in
   # terminal #1.
   # udp://2021 means any IP in this machine.
   go run main/server/smartlog-server.go udp://:2021 file://stdout tcp://localhost:2022
   ```

1. In a third terminal, run:

   ```sh
   # Terminal #3
   # Use the test client to generate some noise and send it over UDP to port 2021.
   go run main/testclient/testclient.go udp://localhost:2021
   ```

After this, you should:

- See the sent messages in terminal #2 (because of the client `file://stdout`)
- See the same messages in `/tmp/out/txt` (because of the client `file:///tmp/out.txt`)
- See the same messages when you point your browser to `http://localhost:8080`.

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

### Stored messages in HTTP clients

HTTP clients store a limited number of messages. The oldest ones are discarded when new messages arrive and the limit is reached. The limit value is the variable `KeepMessages` in the package "smartlog/client/http". To change this value:

```go
import (
  "smartlog/client/http"
)
...
http.KeepMessages = 10000 // store a lot of messages
```

It should be noted that if you need to do this, then maybe you should not log just to an HTTP client, but in parallel also to a different kind - maybe a file client that's not limited by resources (other than diskspace, which is cheap). This can be achieved by:

- Instantiating a forwarding client over TCP or UDP,
- Having a Smartlog server aaccept these messages
- Configuring it to fan out the messages to both an HTTP and a file client.
