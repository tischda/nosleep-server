package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

const PROG_NAME string = "nosleep-server"
const DEFAULT_PORT = 9001

var version string

var flagHelp = flag.Bool("help", false, "displays this help message")
var flagPort = flag.Int("port", DEFAULT_PORT, "RPC server listening port")
var flagDisplay = flag.Bool("display", false, "Force display to stay on")
var flagVersion = flag.Bool("version", false, "print version and exit")

func init() {
	flag.BoolVar(flagHelp, "h", false, "")
	flag.IntVar(flagPort, "p", DEFAULT_PORT, "")
	flag.BoolVar(flagDisplay, "d", false, "")
	flag.BoolVar(flagVersion, "v", false, "")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+PROG_NAME+` [--port <port>] [--display]  | --version | --help

Sets ThreadExecutionState to (ES_CONTINUOUS | ES_SYSTEM_REQUIRED) and
starts an RPC server on 127.0.0.1:`+fmt.Sprintf("%d", DEFAULT_PORT)+`.

You can manage the server using RPC calls to control thread execution states
where possible methods are: Clear, Display, System, Critical, Read and Shutdown.

OPTIONS:

  -d, -display
        Force display to stay on
  -h, -help
        displays this help message
  -p, -port int
        RPC server listening port (default 9001)
  -v, -version
        print version and exit

EXAMPLES:`)

		fmt.Fprintln(os.Stderr, "\n  "+PROG_NAME+` --port 9015 --display

  will set ThreadExecutionState to (ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED)
  and start an RPC server listening on 127.0.0.1:9015.`)
	}
	flag.Parse()

	if flag.Arg(0) == "version" || *flagVersion {
		fmt.Printf("%s version %s\n", PROG_NAME, version)
		return
	}

	if *flagHelp {
		flag.Usage()
		return
	}

	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}

	// register RPC with state manager
	shutdownCh := make(chan struct{})
	manager := &ExecStateManager{rpcShutdownCh: shutdownCh}
	manager.Start()
	defer manager.Stop()

	// set the initial sleep mode
	if *flagDisplay {
		manager.Display(ExecStateRequest{}, &ExecStateReply{})
	} else {
		manager.System(ExecStateRequest{}, &ExecStateReply{})
	}

	// Register RPC server with ExecStateManager methods
	rpc.Register(manager)

	// Configure listener
	address := fmt.Sprintf("127.0.0.1:%d", *flagPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	log.Printf("Nosleep RPC server listening on %s", address)

	// Accept connections until shutdown is triggered
	go func() {
		rpc.Accept(listener)
	}()

	<-shutdownCh
	listener.Close()

	log.Println("Server shutdown complete.")
}
