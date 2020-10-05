package decode

import (
	"encoding/json"
	"fmt"

	"github.com/zencoder/go-smile/domain"
)

func Decode(smile []byte) (string, error) {
	header, err := domain.DecodeHeader(smile)
	if err != nil {
		return "", err
	}

	_, b, err := decodeBytes(smile[header.SizeBytes:])
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(b)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func decodeBytes(smileBytes []byte) ([]byte, interface{}, error) {
	var token = smileBytes[0]
	var tokenClass = token >> 5
	switch tokenClass {
	case 4:
		// Tiny Unicode (2 - 33 bytes; <= 33 characters)
		return readTinyUTF8(smileBytes)
	case 5:
		// Short Unicode (34 - 64 bytes; <= 64 characters)
		return readShortUTF8(smileBytes)
	case 7:
		// Binary / Long text / structure markers (0xF0 - 0xF7 is unused, reserved for future use -- but note, used in key mode)
		return parseBinaryLongTextStructureValues(smileBytes)
	}

	return []byte{}, "", fmt.Errorf("unrecognised token: %X", token)
}

func parseBinaryLongTextStructureValues(smileBytes []byte) ([]byte, interface{}, error) {
	nextByte := smileBytes[0]
	if nextByte == START_OBJECT {
		// Move forward past the START_OBJECT token
		smileBytes = smileBytes[1:]

		// Create an abstract representation of this object
		var object = map[string]interface{}{}

		for smileBytes[0] != END_OBJECT {
			var key, value interface{}
			var err error

			smileBytes, key, err = parseKey(smileBytes)
			if err != nil {
				return smileBytes, object, err
			}

			smileBytes, value, err = decodeBytes(smileBytes)
			if err != nil {
				return smileBytes, object, err
			}

			object[key.(string)] = value
		}
		return smileBytes, object, nil
	}

	return nil, nil, fmt.Errorf("Unknown Byte '%X' in parseBinaryLongTextStructureValues, ignoring...\n", nextByte)
}
