// Package data provides all the methods needed for collecting and saving node
// data to disk.
package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/m-lab/nodeinfo/api"
	"github.com/m-lab/nodeinfo/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Gatherer holds all the information needed about a single data-producing command.
type Gatherer struct {
	Name string
	Cmd  []string
}

// Gather runs the command and gathers the data into the file in the directory.
func (g Gatherer) Gather(crashOnError bool, nodeinfo *api.NodeInfoV1) {
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
	g.gather(nodeinfo)
}

// Save marshals the gathered data, writes it to a file, and returns
// the filename and/or error (if any).
func Save(datadir, datatype string, nodeinfo api.NodeInfoV1) (string, error) {
	b, err := json.Marshal(nodeinfo)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data (error: %v)", err)
	}
	nowUTC := time.Now().UTC()
	dir := fmt.Sprintf("%s/%s/%s", datadir, datatype, nowUTC.Format("2006/01/02"))
	if err := os.MkdirAll(dir, 0o775); err != nil {
		return "", fmt.Errorf("failed to create directory (error: %v)", err)
	}
	file := fmt.Sprintf("%s/%s.json", dir, nowUTC.Format("20060102T150405.000000Z"))
	log.Print(file)
	if err := os.WriteFile(file, b, 0o666); err != nil {
		return file, fmt.Errorf("failed to write file (error: %v)", err)
	}
	return file, nil
}

// gather runs the command. Gather sets up all monitoring, metrics, and
// recovery code, and then gather() does the work.
func (g Gatherer) gather(nodeinfo *api.NodeInfoV1) {
	cmd := api.CmdOut{
		Name:        g.Name,
		CommandLine: strings.Join(g.Cmd, " "),
	}
	log.Printf("   %v\n", cmd.CommandLine)
	out, err := exec.Command(g.Cmd[0], g.Cmd[1:]...).Output()
	if err != nil {
		log.Panicf("failed to run %v (error: %v)", cmd.CommandLine, err)
	}
	cmd.Output = strings.TrimSuffix(string(out), "\n")
	nodeinfo.Commands = append(nodeinfo.Commands, cmd)
}
