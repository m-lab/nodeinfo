// nodeinfo is a utility program for gathering all hw/sw/config data from a
// node that may be operationally relevant.  It is intended to produce lots of
// small files, each with the output of "ifconfig" or "lshw" or another command
// like that. The hope is that by doing this, we will be able to track over
// time what hardware was installed, what software versions were running, and
// how the network was configured on every node in the M-Lab fleet.
//
// nodeinfo reads the list of commands and datatypes in from a config file. It
// rereads the config file every time it runs, to allow that file to be deployed
// as a ConfigMap in kubernetes.
package main

import (
	"context"
	"flag"
	"log"
	"path"
	"time"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/memoryless"
	"github.com/m-lab/go/prometheusx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/go/uniformnames"
	"github.com/m-lab/nodeinfo/config"
	"github.com/m-lab/nodeinfo/metrics"
)

// Command-line flags
var (
	datadir    = flag.String("datadir", "/var/spool/nodeinfo", "The root directory in which to put all produced data")
	once       = flag.Bool("once", false, "Only gather data once")
	smoketest  = flag.Bool("smoketest", false, "Gather every type of data once. Used to test that all configured data types can be gathered.")
	waittime   = flag.Duration("wait", 1*time.Hour, "How long (in expectation) to wait between runs")
	configFile = flag.String("config", "/etc/nodeinfo/config.json", "The name of the config file to load from disk.")

	// A context and associate cancellation function which, when called, should cause main to exit.
	mainCtx, mainCancel = context.WithCancel(context.Background())

	// Contents of this should be filled in as part of parsing commandline flags.
	gatherers config.Config
)

func init() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.LstdFlags)
}

// Runs every data gatherer.
func gather() {
	err := gatherers.Reload()
	if err != nil {
		metrics.ConfigLoadFailures.Inc()
		log.Println("Could not reload the config. Using old config.")
	}
	t := time.Now()
	for _, g := range gatherers.Gatherers() {
		g.Gather(t, *datadir, *smoketest)
	}
}

func main() {
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Could not parse args from environment")

	rtx.Must(uniformnames.Check(path.Base(*datadir)), "The destination directory does not conform to the M-Lab uniform naming conventions")

	metricSrv := prometheusx.MustServeMetrics()
	defer metricSrv.Shutdown(mainCtx)

	var err error
	gatherers, err = config.Create(*configFile)
	rtx.Must(err, "Could not read config on the first try. Shutting down.")
	rtx.Must(
		memoryless.Run(mainCtx, gather, memoryless.Config{Expected: *waittime, Max: 4 * (*waittime), Once: *once || *smoketest}),
		"Bad time arguments.")
}
