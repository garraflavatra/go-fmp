package fmp

func parseVarUint32(payload []byte) uint32 {
	var length uint32
	n := len(payload)
	if n > 4 {
		n = 4 // clamp to max uint32
	}
	for i := range n {
		length <<= 8
		length |= uint32(payload[i])
	}
	return length
}

func decodeByteSeq(payload []byte) string {
	result := ""
	for i := range payload {
		result += string(payload[i] ^ 0x5A)
	}
	return result
}
