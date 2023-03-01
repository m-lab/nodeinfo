// Package config implements all configuration-related logic for this repo.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/m-lab/go/uniformnames"

	"github.com/m-lab/nodeinfo/data"
	"github.com/m-lab/nodeinfo/metrics"
)

// Config contains the configuration of nodeinfo that is stored in a separate
// config file.
type Config interface {
	Reload() error
	Gatherers() []data.Gatherer
}

// Create creates a new config based on the passed-in file name and contents. If
// the file can't be read or parsed, then this will return a non-nil error.
func Create(filename string) (Config, error) {
	c := &fileconfig{
		filename: filename,
	}
	err := c.Reload()
	return c, err
}

// fileconfig contains the full runtime config of nodeinfo.
type fileconfig struct {
	filename  string
	gatherers []data.Gatherer
}

// Reload the list of gatherers from the original config filename. Returns a
// non-nil error if unsuccessful. The config must be well-formed - either the
// whole file is readable and parseable, or the reload will not be successful
// and the list will not be updated.
func (c *fileconfig) Reload() error {
	metrics.ConfigLoadCount.Inc()
	contents, err := ioutil.ReadFile(c.filename)
	if err != nil {
		log.Printf("failed to read %v: %v\n", c.filename, err)
		return err
	}
	var newGatherers []data.Gatherer
	err = json.Unmarshal(contents, &newGatherers)
	if err != nil {
		log.Printf("failed to parse %q", c.filename)
		return err
	}
	for _, g := range newGatherers {
		if g.Name == "" || len(g.Cmd) == 0 {
			log.Printf("%#v is not a valid gatherer", g)
			return fmt.Errorf("%#v is not a valid gatherer", g)
		}
		if err := uniformnames.Check(g.Name); err != nil {
			return err
		}
	}
	c.gatherers = newGatherers
	metrics.ConfigLoadTime.SetToCurrentTime()
	return nil
}

// Gatherers returns a slice of data gatherers. The backing storage for a given
// slice should be immutable.
func (c *fileconfig) Gatherers() []data.Gatherer {
	return c.gatherers
}
