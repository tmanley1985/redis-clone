# Architecture

This is the architecture I hope to support. At least at this time. Being new to this, I may revise this as I become aware of best practices.

```md
[Client] <--TCP--> [Server]
                     |
                     | (Raw input)
                     v
               [RESP Parser]
                     |
                     | (Parsed command)
                     v
             [Command Factory]
                     |
                     | (Command object)
                     v
               [Command Bus]
                     |
                     | (Queued command)
                     v
           [Command Executor(s)]
                     |
                     | (Execution request)
                     v
               [Data Store]
                     |
                     | (Result)
                     v
           [Response Formatter]
                     |
                     | (Formatted response)
                     v
                 [Server]
                     |
                     | (TCP response)
                     v
                 [Client]
```