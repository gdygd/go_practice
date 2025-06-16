package general

const (
	ED_BIG    = 1
	ED_LITTLE = 2
)

func GetNumber(src []byte, pos int, length int, endian int) int {
	var value int = 0

	switch endian {
	case ED_BIG:
		for idx := 0; idx < length; idx++ {
			value = (value * 256) + int(src[pos+idx])
		}
		break
	case ED_LITTLE:
		for idx := length - 1; idx >= 0; idx-- {
			value = (value * 256) + int(src[pos+idx])
		}
		break
	}

	return value
}

func SetNumber(buf []byte, pos int, value int, length int, endian int) {
	switch endian {
	case ED_BIG:
		for idx := pos + length - 1; idx >= pos; idx-- {
			buf[idx] = byte(value % 256)
			value /= 256
		}
		break
	case ED_LITTLE:
		for idx := pos; idx < pos+length; idx++ {
			buf[idx] = byte(value % 256)
			value /= 256
		}
		break
	}
}

func GenLRC(buf []byte, pos int, lastIdx int) byte {
	var lrc byte = 0
	for idx := pos; idx < lastIdx; idx++ {
		lrc ^= buf[idx]
	}

	return lrc
}
