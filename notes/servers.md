# Servers

I've always taken servers for granted. Admittedly, coming from a background in php, python and node it was an afterthought. I knew there was something going on in a lower level of abstraction so I decided to black box it.

The following is an evolving understanding of the basics of writing servers from scratch with a particular focus on Go. I'm going to try to not judge myself too harshly here.

## Flow

This is a very generic flow for a server setup. What you actually do here is entirely dependent upon what you're trying to implement. Writing a program to implement the LLRP protocol is going to vary drastically from hobbling together a concurrent game server. Nevertheless there are techniques that are shared between the two.

In general, the flow will look something like this:

- Listen on a socket
- Listen for _new_ connections
- Optionally store the connections
  - You may need to do this if communication needs to occur between peers. Examples include game servers, and chat applications.
- Handle each connection
  - This is best done in a goroutine so that you're not blocking new connections
- For each connection read the bytes into a buffer
- Do something with the bytes you have in the buffer
  - As stated above, this depends on what you're implementing
- Communication between client/server can be blocking or async using channels

## Buffers

TODO
