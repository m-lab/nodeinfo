// nodeinfo is a utility program for gathering all hw/sw/config data from a
// node that may be operationally relevant.  It is intended to produce lots of
// small files, each with the output of "ifconfig" or "lshw" or another command
// like that. The hope is that by doing this, we will be able to track over
// time what hardware was installed, what software versions were running, and
// how the network was configured on every node in the M-Lab fleet.
//
// nodeinfo reads the list of commands and datatypes in from a config file and
// then opens a webserver on the port specified in the -reload-address argument.
// If the config file changes, you can cause nodeinfo to reload the config by
// sending an HTTP POST to the '/-/reload' url being served from that address.
// You SHOULD NOT expose the reload-address to the world. This pattern allows us
// to reload the config whenever the configmap changes using the
// jimmidyson/configmap-reload:v0.2.2 image.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/m-lab/go/httpx"

	"github.com/m-lab/go/prometheusx"
	"github.com/m-lab/nodeinfo/gatherers"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/memoryless"
	"github.com/m-lab/go/rtx"
)

var (
	reloadAddr  = flag.String("reload-address", "127.0.0.1:9989", "The address to which we should bind the server which serves the reload URL.")
	datadir     = flag.String("datadir", "/var/spool/nodeinfo", "The root directory in which to put all produced data")
	once        = flag.Bool("once", false, "Only gather data once")
	smoketest   = flag.Bool("smoketest", false, "Gather every type of data once. Used to test that all configured data types can be gathered.")
	waittime    = flag.Duration("wait", 1*time.Hour, "How long (in expectation) to wait between runs")
	configFile  = flag.String("config", "/etc/nodeinfo/config.json", "The name of the config file to load from disk.")
	ctx, cancel = context.WithCancel(context.Background())

	// Contents of this should be filled in as part of parsing commandline flags.
	config gatherers.Config
)

func init() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.LstdFlags)
}

// Runs every data gatherer.
func gather() {
	t := time.Now()
	for _, g := range config.Gatherers() {
		g.Gather(t, *datadir, *smoketest)
	}
}

func reloadConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	config.MustReloadConfig()
	w.WriteHeader(200)
}

func main() {
	flag.Parse()
	rtx.Must(flagx.ArgsFromEnv(flag.CommandLine), "Could not parse args from environment")

	defer cancel()

	metricSrv := prometheusx.MustServeMetrics()
	defer metricSrv.Shutdown(ctx)

	reloadHandler := http.NewServeMux()
	reloadHandler.HandleFunc("/-/reload", reloadConfig)
	reloadSrv := &http.Server{
		Handler: reloadHandler,
		Addr:    *reloadAddr,
	}
	httpx.ListenAndServeAsync(reloadSrv)
	defer reloadSrv.Shutdown(ctx)

	config = gatherers.MustCreate(*configFile)

	rtx.Must(
		memoryless.Run(ctx, gather, memoryless.Config{Expected: *waittime, Max: 4 * (*waittime), Once: *once || *smoketest}),
		"Bad time arguments.")
}
