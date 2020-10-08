package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zencoder/go-smile/decode"
	"github.com/zencoder/go-smile/test/testdata"
)

func TestDecode(t *testing.T) {
	filenames, err := testdata.TestFilenames()
	require.NoError(t, err)

	for _, f := range filenames {
		f := f
		t.Run(filepath.Base(f), func(t *testing.T) {
			jsonFile := testdata.LoadTestFile(t, f+".json")
			smileFile := testdata.LoadTestFile(t, f+".smile")

			actualJSON, err := decode.Decode(smileFile)
			require.NoError(t, err, "Error while decoding %q", f)

			require.JSONEq(t, string(jsonFile), actualJSON, "Decoding %q didn't produce the expected result", f)

		})
	}
}
