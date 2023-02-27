package data

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/nodeinfo/api"
)

// Tests are in package data to allow saving data somewhere besides /var/spool/nodeinfo
func TestGather(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestGather")
	rtx.Must(err, "failed to create tempdir")
	defer os.RemoveAll(dir)
	ts := time.Date(2018, 12, 13, 11, 45, 23, 0, time.UTC)
	g := Gatherer{
		Name: "test",
		Cmd:  []string{"echo", "hi"},
	}
	nodeinfo := &api.NodeInfoV1{}
	g.Gather(ts, dir, true, nodeinfo)
	if len(nodeinfo.CommandOutput) != 1 {
		t.Errorf("len(nodeinfo.CommandOutput) = %v, expected 1", len(nodeinfo.CommandOutput))
	}
	co := nodeinfo.CommandOutput[0]
	if co.Name != "test" || co.CommandLine != "echo hi" || co.Output != "hi" {
		t.Errorf("co=%#v, wanted {Name:\"test\", CommandLine:\"echo hi\", Output:\"hi\"}", co)
	}
}

func TestGatherWontCrashWhenItShouldnt(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestGatherWontCrashWhenItShouldnt")
	rtx.Must(err, "failed to create tempdir")
	defer os.RemoveAll(dir)
	g := Gatherer{
		Name: "true",
		Cmd:  []string{"true"},
	}
	g.Gather(time.Now().UTC(), dir, false, &api.NodeInfoV1{})
	// No panic == success
}

func TestGatherWillCrashWhenItShould(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestGatherWillCrashWhenItShould")
	rtx.Must(err, "failed to create tempdir")
	defer os.RemoveAll(dir)
	g := Gatherer{
		Name: "false",
		Cmd:  []string{"false"},
	}
	saveLogFatalf := logFatalf
	logFatalf = func(format string, v ...any) { panic("logFatalf") }
	defer func() {
		r := recover()
		if r == nil {
			t.Error("recover() = nil, expected panic")
		}
		logFatalf = saveLogFatalf
	}()
	g.Gather(time.Now().UTC(), dir, true, &api.NodeInfoV1{})
	// panic == success
}
