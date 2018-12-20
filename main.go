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
	"flag"
	"math"
	"math/rand"
	"time"

	"github.com/m-lab/nodeinfo/data"
)

var (
	once = flag.Bool("once", true, "Only gather data once")
)

// Runs every data gatherer.
func gather() {
	t := time.Now()
	for _, g := range []data.Gatherer{
		{
			Datatype: "lshw",
			Filename: "lshw.json",
			Cmd:      []string{"lshw", "-json"},
		},
		{
			Datatype: "lspci",
			Filename: "lspci.txt",
			Cmd:      []string{"lspci", "-mm", "-vv", "-k", "-nn"},
		},
		{
			Datatype: "ifconfig",
			Filename: "ifconfig.txt",
			Cmd:      []string{"ifconfig", "-a"},
		},
		{
			Datatype: "route",
			Filename: "route-ipv4.txt",
			Cmd:      []string{"route", "-n", "-A", "inet"},
		},
		{
			Datatype: "route",
			Filename: "route-ipv6.txt",
			Cmd:      []string{"route", "-n", "-A", "inet6"},
		},
		{
			Datatype: "uname",
			Filename: "uname.txt",
			Cmd:      []string{"uname", "-a"},
		},
	} {
		g.Gather(t)
	}
}

func main() {
	if *once {
		gather()
	} else {
		for {
			gather()
			time.Sleep(time.Duration(math.Min(rand.ExpFloat64(), 4) * float64(time.Hour)))
		}
	}
}
