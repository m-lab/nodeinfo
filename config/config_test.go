package config_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/m-lab/nodeinfo/config"
	"github.com/m-lab/nodeinfo/data"

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
	c, err := config.Create(dir + "/config.json")
	rtx.Must(err, "Could not read config.json")
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
	rtx.Must(c.Reload(), "Could not reload config")
	g = c.Gatherers()
	if !reflect.DeepEqual(g, expected2) {
		t.Errorf("%v != %v", g, expected2)
	}
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte("bad content"), 0666), "Could not write replacement config")
	if c.Reload() == nil {
		t.Error("We should not have been able to reload the config")
	}
	g = c.Gatherers()
	if !reflect.DeepEqual(g, expected2) {
		t.Errorf("%v != %v", g, expected2)
	}

	incompleteFileContents := []string{
		`[
			{
				"Dataype": "ls",
				"Filename": "ls.txt",
				"Cmd": ["ls", "-l"]
			}
		]
		`,
		`[
			{
				"Datatype": "ls",
				"Filenam": "ls.txt",
				"Cmd": ["ls", "-l"]
			}
		]
		`,
		`[
			{
				"Datatype": "ls",
				"Filename": "ls.txt",
				"Cmb": ["ls", "-l"]
			}
		]
		`,
	}
	for _, inc := range incompleteFileContents {
		rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(inc), 0666), "Could not write replacement config")
		if c.Reload() == nil {
			t.Error("We should not have been able to reload the config")
		}
	}
}

func TestConfigOnBadFile(t *testing.T) {
	_, err := config.Create("/this/file/does/not/exist")
	if err == nil {
		t.Error("This should not have succeeded")
	}
}
