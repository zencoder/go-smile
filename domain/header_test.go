package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeHeaderWithoutData(t *testing.T) {
	_, err := DecodeHeader([]byte{})
	require.EqualError(t, err, "smile format must begin with a 4-byte header")
}

func TestDecodeHeaderWithoutSmilieFails(t *testing.T) {
	_, err := DecodeHeader([]byte{':', '(', '\n', '1'})
	require.EqualError(t, err, "smile format must begin with the ':)' header followed by a newline")
}
