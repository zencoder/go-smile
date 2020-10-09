package decode

import (
	"encoding/json"
	"fmt"

	"github.com/zencoder/go-smile/domain"
)

func DecodeToJSON(smile []byte) (string, error) {
	obj, err := DecodeToObject(smile)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func DecodeToObject(smile []byte) (interface{}, error) {
	header, err := domain.DecodeHeader(smile)
	if err != nil {
		return "", err
	}

	var d decoder
	_, b, err := d.decodeBytes(smile[header.SizeBytes:])
	return b, err
}

type decoder struct {
	sharedState SharedState
}

func (d *decoder) decodeBytes(smileBytes []byte) ([]byte, interface{}, error) {
	var token = smileBytes[0]
	var tokenClass = token >> 5
	switch tokenClass {
	case 0:
		// Short Shared Value String reference
		value, err := d.sharedState.GetSharedValue(int(token&0x1f) - 1)
		return smileBytes[1:], value, err
	case 1:
		// Simple literals, numbers
		return readSimpleLiteral(smileBytes)
	case 2:
		// Tiny ASCII (1 - 32 bytes)
		smileBytes, value, err := readTinyAscii(smileBytes)
		if err == nil {
			d.sharedState.AddSharedValue(value)
		}
		return smileBytes, value, err
	case 3:
		// Short ASCII (33 - 64 bytes)
		smileBytes, value, err := readShortAscii(smileBytes)
		if err == nil {
			d.sharedState.AddSharedValue(value)
		}
		return smileBytes, value, err
	case 4:
		// Tiny Unicode (2 - 33 bytes; <= 33 characters)
		smileBytes, value, err := readTinyUTF8(smileBytes)
		if err == nil {
			d.sharedState.AddSharedValue(value)
		}
		return smileBytes, value, err
	case 5:
		// Short Unicode (34 - 64 bytes; <= 64 characters)
		smileBytes, value, err := readShortUTF8(smileBytes)
		if err == nil {
			d.sharedState.AddSharedValue(value)
		}
		return smileBytes, value, err
	case 6:
		// Small integers (single byte)
		return smileBytes[1:], zigzagDecode(int(token & 0x1f)), nil
	case 7:
		// Binary / Long text / structure markers (0xF0 - 0xF7 is unused, reserved for future use -- but note, used in key mode)
		return d.parseBinaryLongTextStructureValues(smileBytes)
	}

	return []byte{}, "", fmt.Errorf("unrecognised token: %X (Token Class %d)", token, tokenClass)
}

func (d *decoder) parseBinaryLongTextStructureValues(smileBytes []byte) ([]byte, interface{}, error) {
	nextByte := smileBytes[0]
	switch nextByte {
	case START_OBJECT:
		// Move forward past the START_OBJECT token
		smileBytes = smileBytes[1:]

		// Create an abstract representation of this object
		var object = map[string]interface{}{}

		for smileBytes[0] != END_OBJECT {
			var key, value interface{}
			var err error

			smileBytes, key, err = d.parseKey(smileBytes)
			if err != nil {
				return smileBytes, object, err
			}

			smileBytes, value, err = d.decodeBytes(smileBytes)
			if err != nil {
				return smileBytes, object, err
			}

			object[key.(string)] = value
		}
		return smileBytes[1:], object, nil
	case LONG_VARIABLE_ASCII:
		return readVariableLengthText(smileBytes)
	case START_ARRAY:
		smileBytes = smileBytes[1:]
		var array []interface{}
		for smileBytes[0] != END_ARRAY {
			var obj interface{}
			var err error
			smileBytes, obj, err = d.decodeBytes(smileBytes)
			if err != nil {
				return smileBytes, obj, err
			}
			array = append(array, obj)
		}
		return smileBytes[1:], array, nil
	case SHARED_STRING_REFERENCE_LONG_1, SHARED_STRING_REFERENCE_LONG_2, SHARED_STRING_REFERENCE_LONG_3, SHARED_STRING_REFERENCE_LONG_4:
		// Long Shared Value String reference
		var ref = (int(smileBytes[0]&0x03) << 8) | (int(smileBytes[1] & 0xFF))
		value, err := d.sharedState.GetSharedValue(ref)
		return smileBytes[2:], value, err
	case LONG_UTF8:
		return readVariableLengthText(smileBytes)
	}

	return nil, nil, fmt.Errorf("unknown byte '%X' in parseBinaryLongTextStructureValues\n", nextByte)
}
