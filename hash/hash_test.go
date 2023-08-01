package hash

import (
	"testing"
)

func TestSHA256File(t *testing.T) {
	sha256, _ := SHA256File("hash/hash.go")
	t.Log("SHA256:" + sha256)
}

func TestSHA256Str(t *testing.T) {
	sha256 := SHA256Str("abcdefg")
	t.Log("SHA256:" + sha256)
}
