package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

const DEFAULT_PORT = 9001

// https://goreleaser.com/cookbooks/using-main.version/
var (
	name    string
	version string
	date    string
	commit  string
)

// flags
type Config struct {
	network string
	address string
	port    int
	display bool
	help    bool
	version bool
}

func initFlags() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.network, "n", "tcp", "")
	flag.StringVar(&cfg.network, "network", "tcp", "Network type (tcp, tcp4, tcp6, unix, etc.)")
	flag.StringVar(&cfg.address, "a", "127.0.0.1", "")
	flag.StringVar(&cfg.address, "address", "127.0.0.1", "Bind address")
	flag.IntVar(&cfg.port, "p", DEFAULT_PORT, "")
	flag.IntVar(&cfg.port, "port", DEFAULT_PORT, "RPC server listening port")
	flag.BoolVar(&cfg.display, "d", false, "")
	flag.BoolVar(&cfg.display, "display", false, "Force display to stay on")
	flag.BoolVar(&cfg.help, "?", false, "")
	flag.BoolVar(&cfg.help, "help", false, "displays this help message")
	flag.BoolVar(&cfg.version, "v", false, "")
	flag.BoolVar(&cfg.version, "version", false, "print version and exit")
	return cfg
}

func main() {
	log.SetFlags(0)
	cfg := initFlags()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+name+` [OPTIONS]

Sets ThreadExecutionState to (ES_CONTINUOUS | ES_SYSTEM_REQUIRED) and
starts an RPC server on ADDRESS:PORT (default: 127.0.0.1:`+fmt.Sprintf("%d", DEFAULT_PORT)+`).

You can manage the server using RPC calls to control thread execution states
where possible methods are: Clear, Display, System, Critical, Read and Shutdown.

OPTIONS:

  -n, --network string
          Network type: tcp, tcp4, tcp6, unix or unixpacket (default "tcp")
  -a, --address string
          Bind address (default 127.0.0.1)
  -p, --port int
          RPC server listening port (default 9001)
  -d, --display
          Force display to stay on
  -?, --help
          displays this help message
  -v, --version
          print version and exit

EXAMPLES:`)

		fmt.Fprintln(os.Stderr, "\n  "+name+` --port 9015 --display

  will set ThreadExecutionState to (ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED)
  and start an RPC server listening on 127.0.0.1:9015.`)
	}
	flag.Parse()

	if flag.Arg(0) == "version" || cfg.version {
		fmt.Printf("%s %s, built on %s (commit: %s)\n", name, version, date, commit)
		return
	}

	if cfg.help {
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
	if cfg.display {
		if err := manager.Display(ExecStateRequest{}, &ExecStateReply{}); err != nil {
			log.Fatalf("Failed to set initial display state: %v", err)
		}
	} else {
		if err := manager.System(ExecStateRequest{}, &ExecStateReply{}); err != nil {
			log.Fatalf("Failed to set initial system state: %v", err)
		}
	}

	// Register RPC server with ExecStateManager methods
	if err := rpc.Register(manager); err != nil {
		log.Fatalf("Failed to register RPC server: %v", err)
	}

	// Configure listener
	address := fmt.Sprintf("%s:%d", cfg.address, cfg.port)
	listener, err := net.Listen(cfg.network, address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}
	log.Printf("Nosleep RPC server listening on %s (%s)", address, cfg.network)

	// Accept connections until shutdown is triggered
	go func() {
		rpc.Accept(listener)
	}()

	<-shutdownCh
	listener.Close() //nolint:errcheck

	log.Println("Server shutdown complete.")
}
