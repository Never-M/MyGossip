package gossiper

import (
	"os"
	"path/filepath"
)

// TearDown Delete the trace of testcase
func TearDown(name string) {
	_ = os.RemoveAll(filepath.Join("/tmp/", name))
}
