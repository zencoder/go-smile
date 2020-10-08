package testdata

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFilenames() ([]string, error) {
	pattern := filepath.Join(getTestdataDir(), "*.smile")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return files, err
	}

	for i, filename := range files {
		files[i] = strings.TrimSuffix(filename, ".smile")
	}
	return files, nil
}

func LoadTestFile(t *testing.T, filepath string) []byte {
	b, err := ioutil.ReadFile(filepath)
	require.NoError(t, err, "Error reading test file %q", filepath)

	return b
}

func getTestdataDir() string {
	_, testdataFile, _, _ := runtime.Caller(0)
	return filepath.Dir(testdataFile)
}