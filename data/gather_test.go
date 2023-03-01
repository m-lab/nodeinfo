package data

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/nodeinfo/api"
)

// Tests are in package data to allow saving data somewhere besides /var/spool/nodeinfo
func TestGather(t *testing.T) {
	g := Gatherer{
		Name: "test",
		Cmd:  []string{"echo", "hi"},
	}
	nodeinfo := &api.NodeInfoV1{}
	g.Gather(true, nodeinfo)
	if len(nodeinfo.Commands) != 1 {
		t.Errorf("len(nodeinfo.Commands) = %v, expected 1", len(nodeinfo.Commands))
	}
	cmd := nodeinfo.Commands[0]
	if cmd.Name != "test" || cmd.CommandLine != "echo hi" || cmd.Output != "hi" {
		t.Errorf("cmd=%#v, wanted {Name:\"test\", CommandLine:\"echo hi\", Output:\"hi\"}", cmd)
	}
}

func TestGatherInvalidCommand(t *testing.T) {
	g := Gatherer{
		Name: "test",
		Cmd:  []string{"/non/existent/command"},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("recover() = nil, expected panic")
		}
	}()

	nodeinfo := &api.NodeInfoV1{}
	g.Gather(false, nodeinfo)
	if len(nodeinfo.Commands) != 0 {
		t.Errorf("len(nodeinfo.Commands) = %v, expected 0", len(nodeinfo.Commands))
	}
	cmd := nodeinfo.Commands[0]
	if cmd.Name != "test" || cmd.CommandLine != "echo hi" || cmd.Output != "hi" {
		t.Errorf("cmd=%#v, wanted {Name:\"test\", CommandLine:\"echo hi\", Output:\"hi\"}", cmd)
	}
}

func TestGatherCommandFailed(t *testing.T) {
	g := Gatherer{
		Name: "false",
		Cmd:  []string{"false"},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("recover() = nil, expected panic")
		}
	}()
	g.Gather(true, &api.NodeInfoV1{})
	// panic == success
}

func TestSave(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestSave")
	rtx.Must(err, "failed to create tempdir")
	defer os.RemoveAll(dir)
	nodeinfo1 := api.NodeInfoV1{
		Commands: []api.CmdOut{
			{
				Name:        "name1",
				CommandLine: "cmdLine1",
				Output:      "output1 line 1\noutput2 line 2",
			},
			{
				Name:        "name2",
				CommandLine: "cmdLine2",
				Output:      "output1 line 1\noutput2 line 2",
			},
		},
	}
	want := `{"commands":[{"Name":"name1","CommandLine":"cmdLine1","Output":"output1 line 1\noutput2 line 2"},{"Name":"name2","CommandLine":"cmdLine2","Output":"output1 line 1\noutput2 line 2"}]}`
	file, err := Save(dir, "nodeinfo1", nodeinfo1)
	if err != nil {
		t.Errorf("Save() = %v, wanted nil", err)
	}
	got, err := os.ReadFile(file)
	if err != nil {
		t.Errorf("os.ReadFile() = %v, wanted nil", err)
	}
	if string(got) != want {
		t.Errorf("os.ReadFile() = %v, wanted %v", got, want)
	}
}
