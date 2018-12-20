package data

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/m-lab/go/rtx"
)

// Tests are in package data to allow saving data somewhere besides /var/spool/nodeinfo
var (
	g = Gatherer{
		Datatype: "test",
		Filename: "testfile.txt",
		Cmd:      []string{"echo", "hi"},
	}
)

func TestMakedir(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestMakedir")
	rtx.Must(err, "Could not create tempdir")
	oldroot := root
	root = dir
	defer func() {
		root = oldroot
		os.RemoveAll(dir)
	}()
	ts := time.Date(2018, 12, 13, 11, 45, 23, 0, time.UTC)
	g.Gather(ts)
	data, err := ioutil.ReadFile(dir + "/test/2018/12/13/20181213T11:45:23.000Z-testfile.txt")
	if err != nil || string(data) != "hi\n" {
		t.Errorf("Bad filename %v or bad data %q", err, string(data))
	}

}
