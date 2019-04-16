package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/m-lab/go/prometheusx"
	"github.com/m-lab/nodeinfo/data"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/osx"
	"github.com/m-lab/go/rtx"
)

func TestMainOnce(t *testing.T) {
	// Set things up
	ctx, cancel = context.WithCancel(context.Background())
	dir, err := ioutil.TempDir("", "TestMain")
	rtx.Must(err, "Could not create tempdir")
	revertDatadir := osx.MustSetenv("DATADIR", dir)
	revertOnce := osx.MustSetenv("ONCE", "true")
	revertSmoketest := osx.MustSetenv("SMOKETEST", "true")
	og := gatherers
	gatherers = map[string]data.Gatherer{
		"uname": {
			Datatype: "uname",
			Filename: "uname.txt",
			Cmd:      []string{"uname", "-a"},
		},
		"ifconfig": {
			Datatype: "ifconfig",
			Filename: "ifconfig.txt",
			Cmd:      []string{"ifconfig", "-a"},
		},
	}
	defer func() {
		revertOnce()
		revertDatadir()
		revertSmoketest()
		os.RemoveAll(dir)
		cancel()
		gatherers = og
	}()
	*prometheusx.ListenAddress = ":0"

	// Run main.
	main()

	// Verify that some files were created.
	filecount := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filecount++
		}
		return nil
	})
	if filecount == 0 {
		t.Errorf("No files were produced when we ran main.")
	}
}

func TestMainMultiple(t *testing.T) {
	// Set things up
	ctx, cancel = context.WithCancel(context.Background())
	dir, err := ioutil.TempDir("", "TestMain")
	rtx.Must(err, "Could not create tempdir")
	defer os.RemoveAll(dir)

	*datadir = dir
	*once = false
	*smoketest = false
	*waittime = time.Duration(1 * time.Millisecond)
	*prometheusx.ListenAddress = ":0"
	datatypes = flagx.StringArray{"uname", "bad datatype shouldn't crash things"}

	// Run main but sleep for .5s to guarantee that the timer will go off on its
	// own multiple times.
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()
	main()

	// Verify that some files were created.
	filecount := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filecount++
		}
		return nil
	})
	// The timer going off multiple times should produce multiple files. Thanks to
	// randomness, we don't know exactly how many times, but it should definitely
	// be more than once.
	if filecount <= 1 {
		t.Errorf("Not enough files were produced when we ran main.")
	}
}
