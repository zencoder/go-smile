package domain

import (
	"errors"
)

type Header struct {
	Version                           int
	RawBinaryPresent                  bool
	SharedStringValueEncodingEnabled  bool
	SharedPropertyNameEncodingEnabled bool
	SizeBytes                         int
}

/*
Header consists of:

    Constant byte #0: 0x3A (ASCII ':')
    Constant byte #1: 0x29 (ASCII ')')
    Constant byte #2: 0x0A (ASCII linefeed, '\n')
    Variable byte #3, consisting of bits:
        Bits 4-7 (4 MSB): 4-bit version number; 0x00 for current version (note: it is possible that some bits may be reused if necessary)
        Bits 3: Reserved
        Bit 2 (mask 0x04) Whether '''raw binary''' (unescaped 8-bit) values may be present in content
        Bit 1 (mask 0x02): Whether '''shared String value''' checking was enabled during encoding -- if header missing, default value of "false" must be assumed for decoding (meaning parser need not store decoded String values for back referencing)
        Bit 0 (mask 0x01): Whether '''shared property name''' checking was enabled during encoding -- if header missing, default value of "true" must be assumed for decoding (meaning parser MUST store seen property names for possible back references)

And basically first 2 bytes form simple smiley and 3rd byte is a (Unix) linefeed: this to make command-line-tool based identification simple: choice of bytes is not significant beyond visual appearance. Fourth byte contains minimal versioning marker and additional configuration bits.
*/
func DecodeHeader(b []byte) (Header, error) {
	if len(b) < 4 {
		return Header{}, errors.New("smile format must begin with a 4-byte header")
	}

	if b[0] != ':' || b[1] != ')' || b[2] != '\n' {
		return Header{}, errors.New("smile format must begin with the ':)' header followed by a newline")
	}

	var flags = b[3]
	return Header{
		Version:                           int(flags >> 4),
		RawBinaryPresent:                  (0x04 & flags) != 0,
		SharedStringValueEncodingEnabled:  (0x02 & flags) != 0,
		SharedPropertyNameEncodingEnabled: (0x01 & flags) != 0,
		SizeBytes:                         4,
	}, nil
}
