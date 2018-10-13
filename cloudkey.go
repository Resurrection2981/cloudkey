package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/jnovack/cloudkey/src/leds"
	"github.com/tabalt/pidfile"

	_ "github.com/jnovack/cloudkey/fonts"
	_ "github.com/jnovack/cloudkey/screens"
)

var myLeds leds.LEDS

var tags = map[string]string{
	"SYSLOG_IDENTIFIER": "cloudkey",
}

// CmdLineOpts structure for the command line options
type CmdLineOpts struct {
	delay   float64
	reset   bool
	demo    bool
	version bool
	pidfile string
}

var opts CmdLineOpts

var (
	// Version supplied by the linker
	Version = "v0.0.0"
	// Revision supplied by the linker
	Revision = "00000000"
	// GoVersion supplied by the runtime
	GoVersion = runtime.Version()
)

func buildInfo() string {
	return fmt.Sprintf("cloudkey version %s git revision %s go version %s", Version, Revision, GoVersion)
}

func main() {
	fmt.Println(buildInfo())
	clear()

	// Build the screens in the background
	go buildNetwork(0)
	go buildSpeedTest(1)

	// Fire up the loader while the screens build
	load()

	// time.Sleep(2000 * time.Millisecond)
	// Start the carousel!
	myLeds.LED("blue").On()
	myLeds.LED("white").Off()
	startFadeCarousel()
	// display(1)
}

func init() {
	flag.Float64Var(&opts.delay, "delay", 7500, "delay in milliseconds between screens")
	flag.BoolVar(&opts.reset, "reset", false, "reset/clear the screen")
	flag.BoolVar(&opts.demo, "demo", false, "use fake data for display only")
	flag.StringVar(&opts.pidfile, "pidfile", "/var/run/zeromon.pid", "pidfile")
	flag.BoolVar(&opts.version, "version", false, "print version and exit")
	flagutil.SetFlagsFromEnv(flag.CommandLine, "CLOUDKEY")

	if opts.version {
		// already printed version
		os.Exit(0)
	}

	pid, _ := pidfile.Create(opts.pidfile)

	// Setup Service
	// https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/
	fmt.Println("Starting cloudkey service")
	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	signal.Notify(sigs)
	// method invoked upon seeing signal
	go func() {
		s := <-sigs
		fmt.Printf("Received signal '%s', shutting down\n", s)
		fmt.Println("Stopping cloudkey service")
		_ = pid.Clear()
		os.Exit(1)
	}()

	myLeds = leds.LEDS{}
	myLeds.LED("blue").Off()
	myLeds.LED("white").On()
}
