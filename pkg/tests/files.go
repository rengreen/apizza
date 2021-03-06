package tests

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	rand uint32
	mu   sync.Mutex
)

// TempFile returns the path to a temporary file that does not exits.
// Tempfile is essentially a random filename generator.
func TempFile() string {
	return randFile(os.TempDir(), "", "")
}

// NamedTempFile gives the option to give a named temporary file with
//
// will return "/path/to/temp_dir/<prefix>random_filename<suffix>"
func NamedTempFile(prefix, suffix string) string {
	return randFile(os.TempDir(), prefix, suffix)
}

// WithTempFile is a test wrapper that accepts a function with the file
// and testing.T as arguments.
func WithTempFile(test func(string, *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		test(TempFile(), t)
	}
}

// TempDir returns a temporary directory.
func TempDir() string {
	dir := randFile(os.TempDir(), "", "")
	if err := os.Mkdir(dir, 0755); err != nil {
		return os.TempDir()
	}
	return dir
}

// MkTempDir will create a temporary directory in your operating system's
// temp directory
func MkTempDir(name string) string {
	dir := randFile(os.TempDir(), name, "")
	if err := os.Mkdir(dir, 0755); err != nil {
		panic("could not create temp directory " + dir)
	}
	return dir
}

// Parts of this function came from the Go standard library io/ioutil/tempfile.go
func randFile(dir string, prefix, suffix string) (fname string) {
	for i := 0; i < 1000; i++ {
		fname = filepath.Join(dir, prefix+nextRandom()+suffix)
		if _, err := os.Stat(fname); !os.IsNotExist(err) {
			continue
		}
		break
	}
	return fname
}

// This function came from the Go standard library io/ioutil/tempfile.go
func nextRandom() string {
	mu.Lock()
	r := rand
	if r == 0 {
		r = uint32(time.Now().UnixNano() + int64(os.Getpid()))
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	mu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}
