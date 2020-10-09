package decode

import (
	"fmt"
)

/*
Key mode tokens are only used within JSON Object values; if so, they alternate between value tokens (first a key token; followed by either single-value value token or multi-token JSON Object/Array value). A single token denotes end of JSON Object value; all the other tokens are used for expressing JSON Object property name.

Most tokens are single byte: exceptions are 2-byte "long shared String" token, and variable-length "long Unicode String" tokens.

Byte ranges are divides in 4 main sections (64 byte values each):

    0x00 - 0x3F: miscellaneous
        0x00 - 0x1F: not used, reserved for future versions
        0x20: Special constant name "" (empty String)
        0x21 - 0x2F: reserved for future use (unused for now to reduce overlap between values)
        0x30 - 0x33: "Long" shared key name reference (2 byte token); 2 LSBs of the first byte are used as 2 MSB of 10-bit reference (up to 1024) values to a shared name: second byte used for 8 LSB.
            Note: combined values of 0 through 64 are reserved, since there is more optimal representation -- encoder is not to produce such "short long" values; decoder should check that these are not encountered. Future format versions may choose to use these for specific use.
        0x34: Long (not-yet-shared) Unicode name. Variable-length String; token byte is followed by 64 or more bytes, followed by end-of-String marker byte.
            Note: encoding of Strings shorter than 56 bytes should NOT be done using this type: if such sequence is detected it MAY be considered an error. Further, for ASCII names, Strings with length of 56-64 should also use short String notation
        0x35 - 0x39: not used, reserved for future versions
        0x3A: Not used; would be part of header sequence (which is NOT allowed in key mode!)
        0x3B - 0x3F: not used, reserved for future versions
    0x40 - 0x7F: "Short" shared key name reference; names 0 through 63.
    0x80 - 0xBF: Short ASCII names
        0x80 - 0xBF: names consisting of 1 - 64 bytes, all of which represent UTF-8 Ascii characters (MSB not set) -- special case to potentially allow faster decoding
    0xC0 - 0xF7: Short Unicode names
        0xC0 - 0xF7: names consisting of 2 - 57 bytes that can potentially contain UTF-8 multi-byte sequences: encoders are NOT required to guarantee there is one, but for decoding efficiency reasons are recommended to check (that is: decoders on many platforms will be able to handle ASCII-sequences more efficiently than general UTF-8 names)
    0xF8 - 0xFA: reserved (avoid overlap with START/END_ARRAY, START_OBJECT)
    0xFB: END_OBJECT marker
    0xFC - 0xFF: reserved for framing, not used in key mode (used in value mode)

*/
func (d *decoder) parseKey(smileBytes []byte) ([]byte, interface{}, error) {
	nextByte := smileBytes[0]

	if nextByte == EMPTY_STRING {
		return smileBytes, "", nil
	}
	if nextByte >= 0x30 && nextByte <= 0x33 {
		return d.readLongSharedKey(smileBytes)
	}
	if nextByte == 0x34 {
		return readVariableLengthText(smileBytes[1:])
	}
	if nextByte >= 0x40 && nextByte <= 0x7F {
		return d.readShortSharedKey(smileBytes)
	}
	if nextByte >= 0x80 && nextByte <= 0xBF {
		smileBytes, keyName, err := readTinyAscii(smileBytes)
		if err == nil {
			d.sharedState.AddSharedKey(keyName)
		}
		return smileBytes, keyName, err
	}
	if nextByte >= 0xc0 && nextByte <= 0xf7 {
		smileBytes, keyName, err := readShortUTF8Key(smileBytes)
		if err == nil {
			d.sharedState.AddSharedKey(keyName)
		}
		return smileBytes, keyName, err
	}

	return nil, nil, fmt.Errorf("unexpected key token: %X", nextByte)
}

func readShortUTF8Key(smileBytes []byte) ([]byte, interface{}, error) {
	var length = int(smileBytes[0]&0x1F) + 2
	smileBytes = smileBytes[1:]
	return smileBytes[length:], string(smileBytes[:length]), nil
}

func (d *decoder) readLongSharedKey(smileBytes []byte) ([]byte, interface{}, error) {
	var ref = (int(smileBytes[0]&0x03) << 8) | int(smileBytes[1])
	key, err := d.sharedState.GetSharedKey(ref)
	return smileBytes[2:], key, err
}

func (d *decoder) readShortSharedKey(smileBytes []byte) ([]byte, interface{}, error) {
	var ref = int(smileBytes[0] & 0x3f)
	key, err := d.sharedState.GetSharedKey(ref)
	return smileBytes[1:], key, err
}
