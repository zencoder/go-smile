package decode

func readLongUTF8(smileBytes []byte) ([]byte, interface{}, error) {
	var length = 0
	for smileBytes[length] != STRING_END {
		length++
	}

	var s = string(smileBytes[1 : length+1])

	return smileBytes[length+1:], s, nil
}

func readAscii(smileBytes []byte) ([]byte, interface{}, error) {
	var length = int(smileBytes[0]&0x1F) + 1
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
