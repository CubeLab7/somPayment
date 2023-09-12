package somPayment

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"regexp"
)

func decryptAES(data, key []byte) ([]byte, error) {
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

func basicAuth(login, pass string) (basic string) {
	basic = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", login, pass)))
	return
}

func cleanJSONString(input []byte) string {

	// Создаем регулярное выражение, которое находит все управляющие символы и символы, не входящие в диапазон ASCII.
	re := regexp.MustCompile("[[:cntrl:]]|[^\x20-\x7E]")

	// Удаляем все найденные символы.
	cleaned := re.ReplaceAllString(string(input), "")

	return cleaned
}
