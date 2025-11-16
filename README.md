[![Build Status](https://github.com/tischda/nosleep-server/actions/workflows/build.yml/badge.svg)](https://github.com/tischda/nosleep-server/actions/workflows/build.yml)
[![Test Status](https://github.com/tischda/nosleep-server/actions/workflows/test.yml/badge.svg)](https://github.com/tischda/nosleep-server/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/tischda/nosleep-server/badge.svg)](https://coveralls.io/r/tischda/nosleep-server)
[![Linter Status](https://github.com/tischda/nosleep-server/actions/workflows/linter.yml/badge.svg)](https://github.com/tischda/nosleep-server/actions/workflows/linter.yml)
[![License](https://img.shields.io/github/license/tischda/nosleep-server)](/LICENSE)
[![Release](https://img.shields.io/github/release/tischda/nosleep-server.svg)](https://github.com/tischda/nosleep-server/releases/latest)


# nosleep-server

Windows CLI utility (server) that prevents the computer from entering sleep.

The server will prevent the computer from going to sleep by setting
`SetThreadExecutionState`. The client will communication via RPC with
the server to change the sleep mode or shutdown the server.

The main use case is to prevent sleep during a long running task:

1. start nosleep-server in the background
2. run task (eg. backup script)
3. nosleep-client calls server with shutdown request

It's important to note that `SetThreadExecutionState` only applies to the
current thread, so this server runs an `ExecStateManager` that is locked to
a single OS thread. The RPC server uses this `ExecStateManager` to ensure
consistent state accross calls.

## Install

~~~
go install github.com/tischda/nosleep-server@latest
~~~

## Usage

~~~
Usage: nosleep-server [--port <port>] [--display]

Sets ThreadExecutionState to (ES_CONTINUOUS | ES_SYSTEM_REQUIRED) and
starts an RPC server on 127.0.0.1:9001.

You can manage the server using RPC calls to control thread execution states
where possible methods are: Clear, Display, System, Critical, Read and Shutdown.

OPTIONS:

  -p, --port int
        RPC server listening port (default 9001)
  -d, --display
        Force display to stay on
  -?, --help
        displays this help message
  -v, --version
        print version and exit
~~~

## Examples

~~~
nosleep-server --port 9015 --display
~~~

will set ThreadExecutionState to (ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED)
and start an RPC server listening on 127.0.0.1:9015.

You can test the result like this (requires admin rights):

~~~
❯ powercfg -requests
DISPLAY:
None.

SYSTEM:
[PROCESS] \Device\HarddiskVolume5\src\go\nosleep-server\nosleep-server.exe

AWAYMODE:
None.

EXECUTION:
None.

PERFBOOST:
None.

ACTIVELOCKSCREEN:
None.
~~~

## References

* [tischda/nosleep-client](/tischda/nosleep-client)
* [mhbitarafan/go_wakelock](/mhbitarafan/go_wakelock)
* [brandonherzog/nosleep](/brandonherzog/nosleep)
