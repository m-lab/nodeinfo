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
	"flag"
	"log"
	"strings"
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
	smoketest   = flag.Bool("smoketest", false, "Gather every type of data once. Used to test that all data types can be gathered.")
	waittime    = flag.Duration("wait", 1*time.Hour, "How long (in expectation) to wait between runs")
	datatypes   = flagx.StringArray{}
	ctx, cancel = context.WithCancel(context.Background())

	gatherers = map[string]data.Gatherer{
		"lshw": {
			Datatype: "lshw",
			Filename: "lshw.json",
			Cmd:      []string{"lshw", "-json"},
		},
		"lspci": {
			Datatype: "lspci",
			Filename: "lspci.txt",
			Cmd:      []string{"lspci", "-mm", "-vv", "-k", "-nn"},
		},
		"lsusb": {
			Datatype: "lsusb",
			Filename: "lsusb.txt",
			Cmd:      []string{"lsusb", "-v"},
		},
		"ifconfig": {
			Datatype: "ifconfig",
			Filename: "ifconfig.txt",
			Cmd:      []string{"ifconfig", "-a"},
		},
		"route-v4": {
			Datatype: "route",
			Filename: "route-ipv4.txt",
			Cmd:      []string{"route", "-n", "-A", "inet"},
		},
		"route-v6": {
			Datatype: "route",
			Filename: "route-ipv6.txt",
			Cmd:      []string{"route", "-n", "-A", "inet6"},
		},
		"uname": {
			Datatype: "uname",
			Filename: "uname.txt",
			Cmd:      []string{"uname", "-a"},
		},
	}
)

func possibleTypes() []string {
	datatypes := []string{}
	for datatype := range gatherers {
		datatypes = append(datatypes, datatype)
	}
	return datatypes
}

func init() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.LstdFlags)

	flag.Var(&datatypes, "datatype", "What datatype should be collected. This flag can be used multiple times.  The set of possible datatypes is: {"+strings.Join(possibleTypes(), ", ")+"}")
}

// Runs every data gatherer.
func gather() {
	t := time.Now()
	for _, datatype := range datatypes {
		g, ok := gatherers[datatype]
		if ok {
			g.Gather(t, *datadir, *smoketest)
		} else {
			log.Println("Unknown datatype:", datatype)
		}
	}
}

func main() {
	flag.Parse()
	flagx.ArgsFromEnv(flag.CommandLine)
	if *smoketest {
		*once = true
		datatypes = possibleTypes()
	}

	srv := prometheusx.MustServeMetrics()
	defer srv.Close()

	rtx.Must(
		memoryless.Run(ctx, gather, memoryless.Config{Expected: *waittime, Max: 4 * (*waittime), Once: *once}),
		"Bad time arguments.")
}
