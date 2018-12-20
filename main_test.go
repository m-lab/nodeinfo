package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/m-lab/go/osx"
	"github.com/m-lab/go/rtx"
)

func TestMain(t *testing.T) {
	// Set things up
	dir, err := ioutil.TempDir("", "TestMain")
	rtx.Must(err, "Could not create tempdir")
	defer os.RemoveAll(dir)
	revertDatadir := osx.MustSetenv("DATADIR", dir)
	defer revertDatadir()
	revertOnce := osx.MustSetenv("ONCE", "true")
	defer revertOnce()

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
