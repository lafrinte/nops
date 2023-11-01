package str

import "unsafe"

func ToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
