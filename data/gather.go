// Package data provides all the methods needed for collecting and saving node
// data to disk.
package data

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/m-lab/nodeinfo/api"
	"github.com/m-lab/nodeinfo/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var logFatalf = log.Fatalf

// Gatherer holds all the information needed about a single data-producing command.
type Gatherer struct {
	Name string
	Cmd  []string
}

// Gather runs the command and gathers the data into the file in the directory.
func (g Gatherer) Gather(t time.Time, root string, crashOnError bool, nodeinfo *api.NodeInfoV1) {
	// Optionally recover from errors.
	if !crashOnError {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("failed to run %v (error: %q)\n", g, r)
				metrics.GatherErrors.WithLabelValues(g.Name).Inc()
			}
		}()
	}

	// Report metrics.
	metrics.GatherRuns.WithLabelValues(g.Name).Inc()
	timer := prometheus.NewTimer(metrics.GatherRuntime.WithLabelValues(g.Name))
	defer timer.ObserveDuration()

	// Run the command.
	g.gather(t, root, nodeinfo)
}

// gather runs the command. Gather sets up all monitoring, metrics, and
// recovery code, and then gather() does the work.
func (g Gatherer) gather(t time.Time, root string, nodeinfo *api.NodeInfoV1) {
	co := api.CmdOut{
		Name:        g.Name,
		CommandLine: strings.Join(g.Cmd, " "),
	}
	log.Printf("   %v\n", co.CommandLine)
	out, err := exec.Command(g.Cmd[0], g.Cmd[1:]...).Output()
	if err != nil {
		logFatalf("failed to run command (error: %v)", err)
	}
	co.Output = strings.TrimSuffix(string(out), "\n")
	nodeinfo.CommandOutput = append(nodeinfo.CommandOutput, co)
}
