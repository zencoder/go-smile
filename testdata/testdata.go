package testdata

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"
)

func LoadJSONOrgSample1(t *testing.T) []byte {
	return LoadTestData(t, "json-org-sample1.smile")
}

func LoadJSONOrgSample2(t *testing.T) []byte {
	return LoadTestData(t, "json-org-sample2.smile")
}

func LoadJSONOrgSample3(t *testing.T) []byte {
	return LoadTestData(t, "json-org-sample3.smile")
}

func LoadJSONOrgSample4(t *testing.T) []byte {
	return LoadTestData(t, "json-org-sample4.smile")
}

func LoadUnicodeSample(t *testing.T) []byte {
	return LoadTestData(t, "unicode.smile")
}

func LoadTestData(t *testing.T, filename string) []byte {
	_, testdataFile, _, _ := runtime.Caller(0)
	testdataDir := filepath.Dir(testdataFile)
	filename = filepath.Join(testdataDir, filename)

	b, err := ioutil.ReadFile(filename)
	require.NoError(t, err, "Error reading test file %q", filename)

	return b
}
