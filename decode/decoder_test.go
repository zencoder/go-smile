package decode

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zencoder/go-smile/testdata"
)

func TestItCanDecodeUnicode(t *testing.T) {
	smileBytes := testdata.LoadUnicodeSample(t)
	smileJSON := testdata.LoadUnicodeJSON(t)

	actualJSON, err := Decode(smileBytes)
	require.NoError(t, err)

	require.JSONEq(t, smileJSON, actualJSON)
}
