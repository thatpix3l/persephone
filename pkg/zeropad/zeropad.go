package zeropad

func zeroPadBigEndian(bitSize int, buf []byte) []byte {
	for len(buf) < bitSize/8 {
		buf = append([]byte{0}, buf...)
	}
	return buf
}

func zeroPadLittleEndian(bigSize int, buf []byte) []byte {
	for len(buf) < bigSize/8 {
		buf = append(buf, 0)
	}
	return buf
}

func BigEndian64(buf []byte) []byte {
	return zeroPadBigEndian(64, buf)
}

func BigEndian32(buf []byte) []byte {
	return zeroPadBigEndian(32, buf)
}

func BigEndian16(buf []byte) []byte {
	return zeroPadBigEndian(16, buf)
}

func LittleEndian64(buf []byte) []byte {
	return zeroPadLittleEndian(64, buf)
}

func LittleEndian32(buf []byte) []byte {
	return zeroPadLittleEndian(32, buf)
}

func LittleEndian16(buf []byte) []byte {
	return zeroPadLittleEndian(16, buf)
}
