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
	if len(nodeinfo.Commands) != 1 {
		t.Errorf("len(nodeinfo.Commands) = %v, expected 1", len(nodeinfo.Commands))
	}
	cmd := nodeinfo.Commands[0]
	if cmd.Name != "test" || cmd.CommandLine != "echo hi" || cmd.Output != "hi" {
		t.Errorf("cmd=%#v, wanted {Name:\"test\", CommandLine:\"echo hi\", Output:\"hi\"}", cmd)
	}
}

func TestGatherInvalidCommand(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestGather")
	rtx.Must(err, "failed to create tempdir")
	defer os.RemoveAll(dir)
	ts := time.Date(2018, 12, 13, 11, 45, 23, 0, time.UTC)
	g := Gatherer{
		Name: "test",
		Cmd:  []string{"/non/existent/command"},
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

	nodeinfo := &api.NodeInfoV1{}
	g.Gather(ts, dir, false, nodeinfo)
	if len(nodeinfo.Commands) != 0 {
		t.Errorf("len(nodeinfo.Commands) = %v, expected 0", len(nodeinfo.Commands))
	}
	cmd := nodeinfo.Commands[0]
	if cmd.Name != "test" || cmd.CommandLine != "echo hi" || cmd.Output != "hi" {
		t.Errorf("cmd=%#v, wanted {Name:\"test\", CommandLine:\"echo hi\", Output:\"hi\"}", cmd)
	}
}

func TestGatherCommandFailed(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestGatherCommandFailed")
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
