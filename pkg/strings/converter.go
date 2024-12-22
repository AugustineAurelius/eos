package strings

import "unsafe"

func BytesToString(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

func StringToBytes(string_ string) (bytes []byte) {
	return unsafe.Slice(unsafe.StringData(string_), len(string_))
}
