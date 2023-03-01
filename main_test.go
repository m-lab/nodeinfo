package main

import (
	"context"
	"io/ioutil"
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
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			filecount++
		}
		return nil
	})
	return filecount
}

func TestMainOnce(t *testing.T) {
	// Reset global variables into a known-good start state.
	mainCtx, mainCancel = context.WithCancel(context.Background())
	defer mainCancel()

	dir, err := ioutil.TempDir("", "TestMainOnce")
	rtx.Must(err, "failed to create temp data dir")
	defer os.RemoveAll(dir)
	rtx.Must(os.MkdirAll(dir+"/data", 0o777), "failed to create data subdir")

	config := `[{"Name": "uname", "Cmd": ["uname", "-a"]}]`
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(config), 0o666), "failed to write config")

	*datadir = dir + "/data"
	*configFile = dir + "/config.json"
	*once = true
	*smoketest = true
	*waittime = 3 * time.Second
	*prometheusx.ListenAddress = ":0"

	// Run main.
	main()

	// Verify that some files were created inside uname.
	filecount := countFiles(dir + "/data")
	if filecount == 0 {
		t.Errorf("No files were produced when we ran main.")
	}
}

func TestMainMultipleAndReload(t *testing.T) {
	// Reset global variables into a known-good start state.
	mainCtx, mainCancel = context.WithCancel(context.Background())
	defer mainCancel()

	dir, err := ioutil.TempDir("", "TestMainMultiple")
	rtx.Must(err, "failed to create tempdir")
	defer os.RemoveAll(dir)
	rtx.Must(os.MkdirAll(dir+"/data", 0o777), "failed to create data subdir")

	*datadir = dir + "/data"
	*once = false
	*smoketest = false
	*waittime = time.Duration(1 * time.Millisecond)
	*prometheusx.ListenAddress = ":0"
	*configFile = dir + "/config.json"
	config := `[
		{
			"Name": "uname",
			"Cmd":      ["uname", "-a"]
		},
		{
			"Name": "ifconfig",
			"Cmd":      ["ifconfig", "-a"]
		}
	]
	`
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(config), 0o666), "failed to write config")
	rtx.Must(err, "failed to write config")

	// Run main but sleep for .5s to guarantee that the timer will go off on its
	// own multiple times.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		main()
		wg.Done()
	}()

	start := time.Now().UTC()
	unameokay := false
	ifconfigokay := false
	for time.Now().UTC().Sub(start) < time.Second && !(unameokay && ifconfigokay) {
		unameokay = countFiles(dir+"/data") > 1
		ifconfigokay = countFiles(dir+"/data") > 1
	}
	if !ifconfigokay || !unameokay {
		t.Error("Not enough output was produced in a second")
	}

	newConfig := `[
		{
			"Name": "ls",
			"Cmd":      ["ls"]
		}
	]
	`
	rtx.Must(ioutil.WriteFile(dir+"/config.json", []byte(newConfig), 0o666), "failed to write newConfig")
	time.Sleep(500 * time.Millisecond)
	start = time.Now().UTC()
	lsokay := false
	for time.Now().UTC().Sub(start) < time.Second && !lsokay {
		lsokay = countFiles(dir+"/data") > 1
	}
	if !lsokay {
		t.Errorf("Not enough files were produced with the ls config.")
	}

	os.Remove(dir + "/config.json")
	time.Sleep(100 * time.Millisecond)
	// Make sure that the file disappearing does not cause a crash.

	mainCancel()
	wg.Wait()
}
