// Package data provides all the methods needed for collecting and saving node
// data to disk.
package data

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/m-lab/nodeinfo/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/m-lab/go/rtx"
	pipe "gopkg.in/m-lab/pipe.v3"
)

// Gatherer holds all the information needed about a single data-producing command.
type Gatherer struct {
	Datatype string
	Filename string
	Cmd      []string
}

// filename generates the output filename from the timestamp.
func (g Gatherer) filename(t time.Time) string {
	return t.Format("20060102T15:04:05.000Z-") + g.Filename
}

// makeDirectories creates all the required directories to hold the output filename.
func (g Gatherer) makeDirectories(t time.Time, root string) (string, error) {
	dirname := path.Join(root, g.Datatype, t.Format("2006/01/02"))
	return dirname, os.MkdirAll(dirname, 0775)
}

// Gather runs the command and gathers the data into the file in the directory.
func (g Gatherer) Gather(t time.Time, root string, crashOnError bool) {
	// Optionally recover from errors.
	if !crashOnError {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Failed to run %v (error: %q)\n", g, r)
				metrics.GatherErrors.WithLabelValues(g.Datatype).Inc()
			}
		}()
	}

	// Report metrics.
	metrics.GatherRuns.WithLabelValues(g.Datatype).Inc()
	timer := prometheus.NewTimer(metrics.GatherRuntime.WithLabelValues(g.Datatype))
	defer timer.ObserveDuration()

	// Run the command.
	g.gather(t, root)
}

// gather runs the command. Gather sets up all monitoring, metrics, and
// recovery code, and then gather() does the work.
func (g Gatherer) gather(t time.Time, root string) {
	dir, err := g.makeDirectories(t, root)
	rtx.PanicOnError(err, "Could not make %q", dir)
	outputfile := path.Join(dir, g.filename(t))
	log.Print(outputfile)
	command := pipe.Line(
		pipe.Exec(g.Cmd[0], g.Cmd[1:]...),
		pipe.WriteFile(outputfile, 0666))
	rtx.PanicOnError(pipe.Run(command), "Could not gather %s data", g.Datatype)
}
