# Redis Clone

## Motivation

I've been learning how to implement various protocols on top of tcp while learning go. This is the latest venture towards that goal.

## Setup

You can just run make in the root of the project. That will start the server.

In another terminal, you can type: `nc localhost 5001` to start a session.

Then to test, you put this in that terminal individually, and hitting enter after each line:

```
*3
$3
SET
$5
mykey
$5
Hello
```

That's an example of a SET command, but I haven't figured out how to send the actual request: `*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nHello\r\n` without the damned delimiter not being read as an actual line return.