package util

import "bytes"

func FixUTF8(v string) []byte {
	unicodeReplacement := []byte{0xef, 0xbf, 0xbd}

	return bytes.ToValidUTF8([]byte(v), unicodeReplacement)
}
