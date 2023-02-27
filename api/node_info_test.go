package api

import (
	"testing"
)

// TestV1 fails if there is backwards-incompatible change to NodeInfoV1.
func TestV1(t *testing.T) {
	nodeinfo1 := NodeInfoV1{
		Commands: []CmdOut{
			{
				Name:        "name1",
				CommandLine: "cmdLine1",
				Output:      "",
			},
			{
				Name:        "name2",
				CommandLine: "cmdLine2",
				Output:      "",
			},
		},
	}
	t.Logf("nodeinfo1=%#v\n", nodeinfo1)
}
