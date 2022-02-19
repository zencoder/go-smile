package decode

import (
	"fmt"
	"math"
	"math/big"
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
		return readVarInt(smileBytes[1:], true)
	case FLOAT_32:
		return readFloat32(smileBytes[1:])
	case FLOAT_64:
		return readFloat64(smileBytes[1:])
	case BIG_INT:
		return readBigInt(smileBytes[1:])
	default:
		return smileBytes[1:], nil, fmt.Errorf("error reading simple literal byte '%X'", literalByte)
	}
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

func readVarInt(smileBytes []byte, doZigZagDecode bool) ([]byte, int, error) {
	var varInt, i int
	for i = 0; smileBytes[i]&0x80 == 0; i++ {
		varInt = varInt << 7
		varInt |= int(smileBytes[i])
	}
	varInt = varInt << 6
	varInt |= int(smileBytes[i] & 0x3F)

	if doZigZagDecode {
		varInt = zigzagDecode(varInt)
	}
	return smileBytes[i+1:], varInt, nil
}

// The format of BigInts is unusual - the initial numeric header tells us how many bytes will be in the
// resulting number, but the bytes themselves are in the form 0x0xxxxxxx, with the leading 0 to be ignored.
// We then have to concatenate these seven bits together to "fill up" the output byte array.
//
// e.g: We get an input of 00000001 01010101
// Which then gets turned into an output of 00000011 010101... etc.
//
// There's a special case on the final byte, where we just read as many bits (from MSB) as we need to fill up the final
// output byte, then discard the rest.
func readBigInt(smileBytes []byte) ([]byte, interface{}, error) {
	smileBytes, arrayLength, err := readVarInt(smileBytes, false)
	if err != nil {
		return smileBytes, arrayLength, err
	}
	numBytesToRead := int(math.Ceil(float64(arrayLength*8) / 7))

	var binaryString string
	for i := 0; i < numBytesToRead; i++ {
		binaryString += fmt.Sprintf("%07b", smileBytes[i])

		if i == numBytesToRead-1 {
			trailing := len(binaryString) % 8
			binaryString = binaryString[:len(binaryString)-trailing]
		}
	}

	var n = new(big.Int)
	n, ok := n.SetString(binaryString, 2)
	if !ok {
		return nil, nil, fmt.Errorf("SetString: error turning binary string into Big Int.\nBinary String: %s", binaryString)
	}

	return smileBytes[numBytesToRead:], n, nil
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
