// nodeinfo is a utility program for gathering all hw/sw/config data from a
// node that may be operationally relevant.  It is intended to produce lots of
// small files, each with the output of "ifconfig" or "lshw" or another command
// like that. The hope is that by doing this, we will be able to track over
// time what hardware was installed, what software versions were running, and
// how the network was configured on every node in the M-Lab fleet.  Every time
// we turn out to need a new small diagnostic command, that command should be
// added to the list and a new image pushed.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/m-lab/go/prometheusx"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/memoryless"
	"github.com/m-lab/go/rtx"

	"github.com/m-lab/nodeinfo/data"
)

var (
	datadir     = flag.String("datadir", "/var/spool/nodeinfo", "The root directory in which to put all produced data")
	once        = flag.Bool("once", false, "Only gather data once")
	smoketest   = flag.Bool("smoketest", false, "Gather every type of data once. Used to test that all configured data types can be gathered.")
	waittime    = flag.Duration("wait", 1*time.Hour, "How long (in expectation) to wait between runs")
	config      = flagx.FileBytes{}
	ctx, cancel = context.WithCancel(context.Background())

	// gatherers should be filled in by the config-reading process
	gatherers = []data.Gatherer{}
)

func init() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.LstdFlags)
	flag.Var(&config, "config", "The name of the config file containg the json config for nodeinfo.")
}

// Runs every data gatherer.
func gather() {
	log.Println("About to gather", len(gatherers), "times")
	t := time.Now()
	for _, g := range gatherers {
		g.Gather(t, *datadir, *smoketest)
	}
}

func main() {
	flag.Parse()
	flagx.ArgsFromEnv(flag.CommandLine)

	rtx.Must(json.Unmarshal(config, &gatherers), "Could not read config")

	srv := prometheusx.MustServeMetrics()
	defer srv.Close()

	rtx.Must(
		memoryless.Run(ctx, gather, memoryless.Config{Expected: *waittime, Max: 4 * (*waittime), Once: *once || *smoketest}),
		"Bad time arguments.")
}
