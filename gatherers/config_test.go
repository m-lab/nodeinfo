package gatherers_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/m-lab/nodeinfo/data"
	"github.com/m-lab/nodeinfo/gatherers"

	"github.com/m-lab/go/rtx"
)

func TestConfigCreationAndReload(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestConfigCreation")
	rtx.Must(err, "Could not create tempdir")
	defer os.RemoveAll(dir)

	filecontents := `[
		{
			"Datatype": "uname",
			"Filename": "uname.txt",
			"Cmd": ["uname", "-a"]
		},
		{
			"Datatype": "ifconfig",
			"Filename": "ifconfig.txt",
			"Cmd": ["ifconfig"]
		}
	]
	`
	expected := []data.Gatherer{
		{Datatype: "uname", Filename: "uname.txt", Cmd: []string{"uname", "-a"}},
		{Datatype: "ifconfig", Filename: "ifconfig.txt", Cmd: []string{"ifconfig"}},
	}
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(filecontents), 0666), "Could not write config")
	c := gatherers.MustCreate(dir + "/config.json")
	g := c.Gatherers()
	if !reflect.DeepEqual(g, expected) {
		t.Errorf("%v != %v", g, expected)
	}

	filecontents2 := `[
		{
			"Datatype": "ls",
			"Filename": "ls.txt",
			"Cmd": ["ls", "-l"]
		}
	]
	`
	expected2 := []data.Gatherer{
		{Datatype: "ls", Filename: "ls.txt", Cmd: []string{"ls", "-l"}},
	}
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(filecontents2), 0666), "Could not write replacement config")
	c.MustReloadConfig()
	g = c.Gatherers()
	if !reflect.DeepEqual(g, expected2) {
		t.Errorf("%v != %v", g, expected2)
	}
}
