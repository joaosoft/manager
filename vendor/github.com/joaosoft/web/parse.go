package web

func parseHexUint(v []byte) (n uint64, err error) {
	for i, b := range v {
		switch {
		case '0' <= b && b <= '9':
			b = b - '0'
		case 'a' <= b && b <= 'f':
			b = b - 'a' + 10
		case 'A' <= b && b <= 'F':
			b = b - 'A' + 10
		default:
			return 0, ErrorInvalidChunk
		}
		if i == 16 {
			return 0, ErrorInvalidChunk
		}
		n <<= 4
		n |= uint64(b)
	}
	return
}