package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/m-lab/go/prometheusx"

	"github.com/m-lab/go/rtx"
)

func countFiles(dir string) int {
	filecount := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filecount++
		}
		return nil
	})
	return filecount
}

func TestMainOnce(t *testing.T) {
	// Set things up
	ctx, cancel = context.WithCancel(context.Background())

	dir, err := ioutil.TempDir("", "TestMainData")
	rtx.Must(err, "Could not create temp data dir")
	defer os.RemoveAll(dir)

	configDir, err := ioutil.TempDir("", "TestMainConfig")
	rtx.Must(err, "Could not create temp config dir")
	defer os.RemoveAll(configDir)

	config := `[{"Datatype": "uname", "Filename": "uname.txt", "Cmd": ["uname", "-a"]}]`
	rtx.Must(ioutil.WriteFile(configDir+"/config.json", []byte(config), 0666), "Could not write config")

	*datadir = dir
	*configFile = configDir + "/config.json"
	*once = true
	*smoketest = true
	*prometheusx.ListenAddress = ":0"
	*reloadAddr = ":0"

	// Run main.
	main()

	// Verify that some files were created inside uname.
	filecount := countFiles(dir + "/uname")
	if filecount == 0 {
		t.Errorf("No files were produced when we ran main.")
	}
}

func TestMainMultipleAndReload(t *testing.T) {
	// Set things up
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	dir, err := ioutil.TempDir("", "TestMainMultiple")
	rtx.Must(err, "Could not create tempdir")
	defer os.RemoveAll(dir)

	*datadir = dir
	*once = false
	*smoketest = false
	*waittime = time.Duration(1 * time.Millisecond)
	*prometheusx.ListenAddress = ":0"
	*reloadAddr = "127.0.0.1:12345"
	*configFile = dir + "/config.json"
	config := `[
		{
			"Datatype": "uname",
			"Filename": "uname.txt",
			"Cmd":      ["uname", "-a"]
		},
		{
			"Datatype": "ifconfig",
			"Filename": "ifconfig.txt",
			"Cmd":      ["ifconfig", "-a"]
		}
	]
	`
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(config), 0666), "Could not write config")
	rtx.Must(err, "Could not write config")

	// Run main but sleep for .5s to guarantee that the timer will go off on its
	// own multiple times.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		main()
		wg.Done()
	}()

	time.Sleep(500 * time.Millisecond)
	filecount := countFiles(dir + "/uname")
	if filecount <= 1 {
		t.Errorf("Not enough files were produced when we ran main.")
	}
	filecount = countFiles(dir + "/ifconfig")
	if filecount <= 1 {
		t.Errorf("Not enough files were produced when we ran main.")
	}

	newConfig := `[
		{
			"Datatype": "ls",
			"Filename": "ls.txt",
			"Cmd":      ["ls"]
		}
	]
	`
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(newConfig), 0666), "Could not write newConfig")
	resp, err := http.Get("http://127.0.0.1:12345/-/reload")
	if err == nil && resp.StatusCode == 200 {
		t.Error("We should not be able to GET that url")
	}
	resp, err = http.Post("http://127.0.0.1:12345/-/reload", "application/json", &bytes.Buffer{})
	if err != nil || resp.StatusCode != 200 {
		t.Error("We should have been able to POST to that url")
	}
	rtx.Must(err, "Could not reload the config")
	time.Sleep(500 * time.Millisecond)
	filecount = countFiles(dir + "/ls")
	if filecount <= 1 {
		t.Errorf("Not enough files were produced when we ran main.")
	}

	cancel()
	wg.Wait()
}
