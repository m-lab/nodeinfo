package gatherers

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/nodeinfo/data"
	"github.com/m-lab/nodeinfo/metrics"
)

// Config contains the configuration of nodeinfo that is stored in a separate
// config file. It is broken out into a threadsafe object here to enable safe
// reloads of the configmap.
//
// Implementations of this interface should be designed to be threadsafe.
type Config interface {
	MustReloadConfig()
	Gatherers() []data.Gatherer
}

// MustCreate creates a new config based on the passed-in file name and
// contents. If the file can't be read or parsed, then this will log.Fatal.
func MustCreate(filename string) Config {
	c := &fileconfig{
		filename: filename,
	}
	c.MustReloadConfig()
	return c
}

// fileconfig contains the full runtime config of nodeinfo.
type fileconfig struct {
	filename  string
	mutex     sync.Mutex
	gatherers []data.Gatherer
}

// MustReloadConfig reloads the list of gatherers from the original config
// filename. Crashes on error.
func (c *fileconfig) MustReloadConfig() {
	metrics.ConfigLoadCount.Inc()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	contents, err := ioutil.ReadFile(c.filename)
	rtx.Must(err, "Could not read config file: %q", c.filename)
	var g []data.Gatherer
	rtx.Must(json.Unmarshal(contents, &g), "Could not parse %q", c.filename)
	c.gatherers = g
	metrics.ConfigLoadTime.SetToCurrentTime()
}

// Gatherers returns a slice of data gatherers. The backing storage for a given
// slice should be immutable. The only thing that might happen is that the slice
// will be reinitialized to point at different backing storage. Therefore, once
// the slice has been copied/returned, it is safe to be used without a lock.
func (c *fileconfig) Gatherers() []data.Gatherer {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.gatherers
}
