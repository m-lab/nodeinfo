package main

import (
	"flag"
	"math"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/m-lab/go/rtx"
	pipe "gopkg.in/m-lab/pipe.v3"
)

var (
	once = flag.Bool("once", true, "Only run the check once")
)

type DataGatherer struct {
	root     string
	datatype string
	filename string
	cmd      []string
}

func (d DataGatherer) Filename(t time.Time) string {
	return t.Format("20060102T15:04:05Z-") + d.filename
}

func (d DataGatherer) MakeDirectories(t time.Time) (string, error) {
	dirname := path.Join(d.root, d.datatype, t.Format("2006/01/02"))
	return dirname, os.MkdirAll(dirname, 0775)
}

func (d DataGatherer) Gather() {
	t := time.Now()
	dir, err := d.MakeDirectories(t)
	rtx.Must(err, "Could not make %q", dir)
	outputfile := path.Join(dir, d.Filename(t))
	command := pipe.Line(
		pipe.Exec(d.cmd[0], d.cmd[1:]...),
		pipe.WriteFile(outputfile, 0666))
	rtx.Must(pipe.Run(command), "Could not gather %s data", d.datatype)
}

func main() {
	for {
		for _, g := range []DataGatherer{
			{
				root:     "/var/spool/hardware",
				datatype: "lshw",
				filename: "lshw.json",
				cmd:      []string{"lshw", "-json"},
			},
			{
				root:     "/var/spool/hardware",
				datatype: "lspci",
				filename: "lspci.txt",
				cmd:      []string{"lspci", "-mm", "-vv", "-k", "-nn"},
			},
			{
				root:     "/var/spool/configuration",
				datatype: "ifconfig",
				filename: "ifconfig.txt",
				cmd:      []string{"ifconfig", "-a"},
			},
			{
				root:     "/var/spool/configuration",
				datatype: "route",
				filename: "route-ipv4.txt",
				cmd:      []string{"route", "-n", "-A", "inet"},
			},
			{
				root:     "/var/spool/configuration",
				datatype: "route",
				filename: "route-ipv6.txt",
				cmd:      []string{"route", "-n", "-A", "inet6"},
			},
		} {
			g.Gather()
		}
		if *once {
			break
		} else {
			time.Sleep(time.Duration(math.Min(rand.ExpFloat64(), 4) * float64(time.Hour)))
		}
	}
}
