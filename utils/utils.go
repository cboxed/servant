package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"os"
)

func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func StringToSha512(str string) string {
	tmp := sha512.Sum512([]byte(str))
	return hex.EncodeToString(tmp[:])
}
