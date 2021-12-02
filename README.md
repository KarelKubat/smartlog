# smartlog

`smartlog` is a package for Go to make setting up logging easier. Log statements can be processed locally (to `stdout` or a file) or sent remotely to a server over TCP or UDP.

`smartlog` contains all support code to embed such logging into your programs:

- To embed client code into some program that can use smart logging,
- To embed server code into some other program for further processing.

This is still very much **work in progress**. This github repository exists just because I don't want to loose work in the case of a local crash. Come back later.

Still, if you want to see what it does, check out the testing server or client under `test/`.