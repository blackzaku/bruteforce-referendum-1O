package o1st

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
)

func auxMD5(prev, password, salt []byte) []byte {
	a := make([]byte, len(prev)+len(password)+len(salt))
	copy(a, prev)
	copy(a[len(prev):], password)
	copy(a[len(prev)+len(password):], salt)
	h := md5.New()
	h.Write(a)
	return h.Sum(nil)
}


func extractOpenSSLCreds(password, salt []byte) ([]byte, []byte) {
	m := make([]byte, 48)
	prev := []byte{}
	for i := 0; i < 3; i++ {
		prev = auxMD5(prev, password, salt)
		copy(m[i*16:], prev)
	}
	return m[:32], m[32:] // key + iv
}


// ciphertext, _ := hex.DecodeString(text)
// password = []byte(password)
func DecryptAES256CBC(ciphertext, password []byte) []byte{
	key, iv := extractOpenSSLCreds([]byte(password), []byte{})
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	// This does it in-place
	mode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext
}