package data_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/nodeinfo/data"
)

// Tests are in package data to allow saving data somewhere besides /var/spool/nodeinfo
func TestGather(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestGather")
	rtx.Must(err, "Could not create tempdir")
	defer os.RemoveAll(dir)
	ts := time.Date(2018, 12, 13, 11, 45, 23, 0, time.UTC)
	g := data.Gatherer{
		Datatype: "test",
		Filename: "testfile.txt",
		Cmd:      []string{"echo", "hi"},
	}
	g.Gather(ts, dir, true)
	data, err := ioutil.ReadFile(dir + "/test/2018/12/13/20181213T11:45:23.000Z-testfile.txt")
	if err != nil || string(data) != "hi\n" {
		t.Errorf("Bad filename %v or bad data %q", err, string(data))
	}
}

func TestGatherWontCrashWhenItShouldnt(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestGatherWontCrashWhenItShouldnt")
	rtx.Must(err, "Could not create tempdir")
	defer os.RemoveAll(dir)
	g := data.Gatherer{
		Datatype: "false",
		Filename: "false.txt",
		Cmd:      []string{"false"},
	}
	g.Gather(time.Now(), dir, false)
	// No panic == success
}

func TestGatherWillCrashWhenItShould(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestGatherWillCrashWhenItShould")
	rtx.Must(err, "Could not create tempdir")
	defer os.RemoveAll(dir)
	g := data.Gatherer{
		Datatype: "false",
		Filename: "false.txt",
		Cmd:      []string{"false"},
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("We should have had a panic here")
		}
	}()
	g.Gather(time.Now(), dir, true)
	// panic == success
}
