package fmp

func addIf(cond bool, val uint64) uint64 {
	if cond {
		return val
	}
	return 0
}

func decodeVarUint64(payload []byte) uint64 {
	var length uint64
	n := min(len(payload), 8) // clamp to uint64
	for i := range n {
		length <<= 8
		length |= uint64(payload[i])
	}
	return length
}

func decodeString(payload []byte) string {
	result := ""
	for i := range payload {
		result += string(payload[i] ^ 0x5A)
	}
	return result
}
