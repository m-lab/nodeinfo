package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

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
// non-nil error if unsuccessful.
func (c *fileconfig) Reload() error {
	metrics.ConfigLoadCount.Inc()
	contents, err := ioutil.ReadFile(c.filename)
	if err != nil {
		log.Println("Could not read file")
		return err
	}
	var newGatherers []data.Gatherer
	err = json.Unmarshal(contents, &newGatherers)
	if err != nil {
		log.Printf("Could not parse %q", c.filename)
		return err
	}
	for _, g := range newGatherers {
		if len(g.Cmd) == 0 || g.Datatype == "" || g.Filename == "" {
			log.Printf("%#v is not a valid gatherer", g)
			return fmt.Errorf("%#v is not a valid gatherer", g)
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
