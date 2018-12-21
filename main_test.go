package main

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/m-lab/nodeinfo/data"

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
	og := gatherers
	gatherers = []data.Gatherer{
		{
			Datatype: "uname",
			Filename: "uname.txt",
			Cmd:      []string{"uname", "-a"},
		},
		{
			Datatype: "ifconfig",
			Filename: "ifconfig.txt",
			Cmd:      []string{"ifconfig", "-a"},
		},
	}
	defer func() {
		revertOnce()
		revertDatadir()
		os.RemoveAll(dir)
		cancel()
		gatherers = og
	}()

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
	revertDatadir := osx.MustSetenv("DATADIR", dir)
	revertOnce := osx.MustSetenv("WAIT", "10ms")
	og := gatherers
	gatherers = []data.Gatherer{
		{
			Datatype: "uname",
			Filename: "uname.txt",
			Cmd:      []string{"uname", "-a"},
		},
		{
			Datatype: "ifconfig",
			Filename: "ifconfig.txt",
			Cmd:      []string{"ifconfig", "-a"},
		},
	}
	defer func() {
		revertOnce()
		revertDatadir()
		os.RemoveAll(dir)
		gatherers = og
	}()

	// Run main but sleep for .5 seconds to guarantee that the timer will go off on
	// its own at least once.
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
	if filecount == 0 {
		t.Errorf("No files were produced when we ran main.")
	}
}
