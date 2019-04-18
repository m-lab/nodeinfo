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

	// Run main.
	main()

	// Verify that some files were created inside uname.
	filecount := countFiles(dir + "/uname")
	if filecount == 0 {
		t.Errorf("No files were produced when we ran main.")
	}
}

func TestMainMultipleAndReload(t *testing.T) {
	// Reset global variables into a known-good start state.
	mainCtx, mainCancel = context.WithCancel(context.Background())
	defer mainCancel()

	dir, err := ioutil.TempDir("", "TestMainMultiple")
	rtx.Must(err, "Could not create tempdir")
	defer os.RemoveAll(dir)

	*datadir = dir
	*once = false
	*smoketest = false
	*waittime = time.Duration(1 * time.Millisecond)
	*prometheusx.ListenAddress = ":0"
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

	start := time.Now()
	unameokay := false
	ifconfigokay := false
	for time.Now().Sub(start) < time.Second && !(unameokay && ifconfigokay) {
		unameokay = countFiles(dir+"/uname") > 1
		ifconfigokay = countFiles(dir+"/ifconfig") > 1
	}
	if !ifconfigokay || !unameokay {
		t.Error("Not enough output was produced in a second")
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
	time.Sleep(500 * time.Millisecond)
	start = time.Now()
	lsokay := false
	for time.Now().Sub(start) < time.Second && !lsokay {
		lsokay = countFiles(dir+"/ls") > 1
	}
	if !lsokay {
		t.Errorf("Not enough files were produced when we ran main.")
	}

	os.Remove(dir + "/config.json")
	time.Sleep(100 * time.Millisecond)
	// Make sure that the file disappearing does not cause a crash.

	mainCancel()
	wg.Wait()
}
