package main

import (
	"flag"
	"fmt"
	"log"
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
	logPath string
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
	flag.StringVar(&cfg.logPath, "l", "", "")
	flag.StringVar(&cfg.logPath, "log", "", "Write logs to a file instead of stdout")
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
where possible commands are: Clear, Display, System, Critical, Read and Shutdown.

Another way to control the server is by registering/unregistering processes.
The server will automatically shut down when the last process is unregistered.

OPTIONS:

  -n, --network string
          Network type: tcp, tcp4, tcp6, unix or unixpacket (default "tcp")
  -a, --address string
          Bind address (default 127.0.0.1)
  -p, --port int
          RPC server listening port (default 9001)
  -d, --display
          Force display to stay on
  -l, --log path
          Write logs to a file instead of stdout
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

	if cfg.logPath != "" {
		logFile, err := os.OpenFile(cfg.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			log.Fatalf("Failed to open log file %s: %v", cfg.logPath, err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
		log.SetFlags(log.LstdFlags)
		log.Printf("----------------- SERVER START [pid=%d] -----------------", os.Getpid())
	}

	log.Printf("%s %s starting...\n", name, version)
	serve(cfg)
}
