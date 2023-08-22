package service

import (
	"crypto/aes"
)

const (
	successParam = "success"
	failParam    = "fail"

	USDCode = 840
)

func DecryptAES(data, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	pKCS5UnPadding(decrypted)

	return decrypted, nil
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func pKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}
