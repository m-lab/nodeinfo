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
	rtx.Must(err, "failed to create tempdir")
	defer os.RemoveAll(dir)

	filecontents := `[
		{
			"Name": "uname",
			"Cmd": ["uname", "-a"]
		},
		{
			"Name": "ifconfig",
			"Cmd": ["ifconfig"]
		}
	]
	`
	expected := []data.Gatherer{
		{Name: "uname", Cmd: []string{"uname", "-a"}},
		{Name: "ifconfig", Cmd: []string{"ifconfig"}},
	}
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(filecontents), 0o666), "failed to write config")
	c, err := config.Create(dir + "/config.json")
	rtx.Must(err, "failed to read config.json")
	g := c.Gatherers()
	if !reflect.DeepEqual(g, expected) {
		t.Errorf("%v != %v", g, expected)
	}

	filecontents2 := `[
		{
			"Name": "ls",
			"Cmd": ["ls", "-l"]
		}
	]
	`
	expected2 := []data.Gatherer{
		{Name: "ls", Cmd: []string{"ls", "-l"}},
	}
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(filecontents2), 0o666), "failed to write replacement config")
	rtx.Must(c.Reload(), "failed to reload config")
	g = c.Gatherers()
	if !reflect.DeepEqual(g, expected2) {
		t.Errorf("%v != %v", g, expected2)
	}
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte("bad content"), 0o666), "failed to write replacement config")
	if c.Reload() == nil {
		t.Error("We should not have been able to reload the config")
	}
	g = c.Gatherers()
	if !reflect.DeepEqual(g, expected2) {
		t.Errorf("%v != %v", g, expected2)
	}

	incompleteFileContents := []string{
		// Mane is not Name
		`[
			{
				"Mane": "ls",
				"Cmd": ["ls", "-l"]
			}
		]
		`,
		// Cmb is not Cmd
		`[
			{
				"Name": "ls",
				"Cmb": ["ls", "-l"]
			}
		]
		`,
		// Name does not conform to the uniform naming conventions.
		`[
			{
				"Name": "ls-a-lot",
				"Cmd": ["ls", "-l"]
			}
		]
		`,
	}
	for _, inc := range incompleteFileContents {
		rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(inc), 0o666), "failed to write replacement config")
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
