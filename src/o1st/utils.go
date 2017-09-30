package o1st

import (
	"crypto/sha256"
	"encoding/hex"
	"bytes"
)

const BUCLE int = 1714


func bucleHash(key string) string {
	var h [256]byte
	b := sha256.Sum256([]byte(key))
	hex.Encode(h[:], b[:])
	for i := 1; i < BUCLE; i++ {
		b = sha256.Sum256(h[:])
		hex.Encode(h[:], b[:])
	}
	return string(h[:])
}

func hashByte(text []byte) [32]byte {
	b := sha256.Sum256(text[:])
	return b
}

func bucleHashByte(key []byte) [32]byte {
	var h [256]byte
	b := sha256.Sum256(key[:])
	hex.Encode(h[:], b[:])
	for i := 1; i < BUCLE; i++ {
		b = sha256.Sum256(h[:])
		hex.Encode(h[:], b[:])
	}
	return b
}
func hashHex(text []byte) (h [64]byte) {
	b := sha256.Sum256(text[:])
	hex.Encode(h[:], b[:])
	return
}

func loopHashHex(key []byte) (h [64]byte) {
	b := sha256.Sum256(key[:])
	hex.Encode(h[:], b[:])
	for i := 1; i < BUCLE; i++ {
		b = sha256.Sum256(h[:])
		hex.Encode(h[:], b[:])
	}
	return
}

func hash(text string) string {
	var h [256]byte
	b := sha256.Sum256([]byte(text))
	hex.Encode(h[:], b[:])
	return string(h[:])
}

func Check(dni string, date string, zip string) string {
	key := dni + date + zip;
	loopSha256 := loopHashHex([]byte(key))
	firstSha256 := hashHex(loopSha256[:])
	secondSha256 := hashByte(firstSha256[:])
	lines := data[int(secondSha256[0]) * 256 + int(secondSha256[1])]
	for _,line := range lines {
		if bytes.Equal(line[0: 30], secondSha256[2: 32]) {
			return string(DecryptAES256CBC(line[30:], firstSha256[:]))
		}
	}
	return ""
}