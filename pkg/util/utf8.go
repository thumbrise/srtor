package util

import "bytes"

func ToUTF8FixedBytes(v string) []byte {
	unicodeReplacement := []byte{0xef, 0xbf, 0xbd}

	return bytes.ToValidUTF8([]byte(v), unicodeReplacement)
}
