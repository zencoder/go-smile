package decode

import (
	"fmt"
	"math"
)

func readVariableLengthText(smileBytes []byte) ([]byte, interface{}, error) {
	var length = 0
	for smileBytes[length] != STRING_END {
		length++
	}

	var s = string(smileBytes[1:length])
	return smileBytes[length+1:], s, nil
}

func readSimpleLiteral(smileBytes []byte) ([]byte, interface{}, error) {
	var literalByte = smileBytes[0]
	switch literalByte {
	case EMPTY_STRING:
		return smileBytes[1:], "", nil
	case NULL:
		return smileBytes[1:], nil, nil
	case FALSE:
		return smileBytes[1:], false, nil
	case TRUE:
		return smileBytes[1:], true, nil
	case INT_32, INT_64:
		return readVarInt(smileBytes[1:])
	case FLOAT_32:
		return readFloat32(smileBytes[1:])
	case FLOAT_64:
		return readFloat64(smileBytes[1:])
	default:
		return smileBytes[1:], nil, fmt.Errorf("error reading simple literal byte '%X'", literalByte)
	}
	// TODO: BigInteger
	// TODO: BigDecimal
}

func readFloat32(smileBytes []byte) ([]byte, interface{}, error) {
	var intBits = uint32(smileBytes[0])
	for i := 1; i < 5; i++ {
		intBits = (intBits << 7) + uint32(smileBytes[i])
	}
	return smileBytes[5:], math.Float32frombits(intBits), nil
}

func readFloat64(smileBytes []byte) ([]byte, interface{}, error) {
	var intBits = uint64(smileBytes[0])
	for i := 1; i < 10; i++ {
		intBits = (intBits << 7) + uint64(smileBytes[i])
	}
	return smileBytes[10:], math.Float64frombits(intBits), nil
}

func readVarInt(smileBytes []byte) ([]byte, interface{}, error) {
	var varInt, i int
	for i = 0; smileBytes[i]&0x80 == 0; i++ {
		varInt = varInt << 7
		varInt |= int(smileBytes[i])
	}
	varInt = varInt << 6
	varInt |= int(smileBytes[i] & 0x3F)

	return smileBytes[i+1:], zigzagDecode(varInt), nil
}

func zigzagDecode(varInt int) int {
	return ((varInt) >> 1) ^ (-((varInt) & 1))
}

func readTinyAscii(smileBytes []byte) ([]byte, interface{}, error) {
	var length = int(smileBytes[0]&0x1F) + 1
	var s = string(smileBytes[1 : length+1])

	return smileBytes[length+1:], s, nil
}

func readShortAscii(smileBytes []byte) ([]byte, interface{}, error) {
	var length = int(smileBytes[0]&0x1F) + 33
	var s = string(smileBytes[1 : length+1])

	return smileBytes[length+1:], s, nil
}

func readTinyUTF8(smileBytes []byte) ([]byte, interface{}, error) {
	var length = int(smileBytes[0]&0x1F) + 2
	var s = string(smileBytes[1 : length+1])

	return smileBytes[length+1:], s, nil
}

func readShortUTF8(smileBytes []byte) ([]byte, interface{}, error) {
	var length = int(smileBytes[0]&0x1F) + 34
	var s = string(smileBytes[1 : length+1])

	return smileBytes[length+1:], s, nil
}
