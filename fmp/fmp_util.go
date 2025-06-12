package fmp

func parseVarUint64(payload []byte) uint64 {
	var length uint64
	n := min(len(payload), 8) // clamp to uint64
	for i := range n {
		length <<= 8
		length |= uint64(payload[i])
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
