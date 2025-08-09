# nosleep-server

Windows CLI utility (server) that prevents the computer from entering sleep.

The server will prevent the computer from going to sleep by setting
`SetThreadExecutionState`. The client will communication via RPC with
the client to change the sleep mode or shutdown the server.

The main use case is to prevent sleep during a long running task:

1. start server in the background
2. run task (eg. backup script)
3. client calls server with shutdown request


### Install

There are no dependencies.

~~~
go install github.com/tischda/nosleep-server@latest
~~~

### Usage

~~~
~~~

### References

[mhbitarafan/go_wakelock](/mhbitarafan/go_wakelock)