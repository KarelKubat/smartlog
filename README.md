# Smartlog

`smartlog` is a package for Go to make setting up logging easier. Log statements can be processed locally (to `stdout` or a file) or sent remotely to a server over TCP or UDP.

`smartlog` contains all support code to embed such logging into your programs:

- To embed client code into some program that can use smart logging,
- To embed server code into some other program for further processing.

This is still very much **work in progress**. This github repository exists just because I don't want to loose work in the case of a local crash. Come back later.

Still, if you want to see what it does, check out the testing server or client under `test/`.

## Smartlog Clients

### Time Stamps

Any client-side invocation like `client.Info("hello world")` leads to a message which has the timestamp. Two settings can be controlled:

- The timestamp format: the default is `"2006-01-02 15:04:05 MST"` (see e.g. https://www.geeksforgeeks.org/time-formatting-in-golang/)
- Whether the time is displayed relative to localtime or to UTC: the default is `false`, the time is displayed relative to localtime.

To change the defaults, simply modify the global variables in `smartlog/msg`:

```go
import (
    "time"
    "smartlog/msg"
)    
// ...
msg.DefaultTimeFormat = time.RFC3339 // format: "2006-01-02T15:04:05Z07:00"
msg.UTCTime = true                   // relative to UTC
```