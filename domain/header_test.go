package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zencoder/go-smile/testdata"
)

func TestDecodeHeader(t *testing.T) {
	b := testdata.LoadJSONOrgSample1(t)

	header, err := DecodeHeader(b)
	require.NoError(t, err)

	require.Equal(
		t,
		Header{
			Version:                           0,
			RawBinaryPresent:                  false,
			SharedStringValueEncodingEnabled:  true,
			SharedPropertyNameEncodingEnabled: true,
			SizeBytes:                         4,
		},
		header,
	)
}

func TestDecodeHeaderWithoutData(t *testing.T) {
	_, err := DecodeHeader([]byte{})
	require.EqualError(t, err, "smile format must begin with a 4-byte header")
}

func TestDecodeHeaderWithoutSmilieFails(t *testing.T) {
	_, err := DecodeHeader([]byte{':', '(', '\n', '1'})
	require.EqualError(t, err, "smile format must begin with the ':)' header followed by a newline")
}
