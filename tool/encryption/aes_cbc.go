package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
)

func AESCBCEncrypt(secretKey SecretKey, origin_data string) string {
	iv := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}
	log.Println(hex.EncodeToString(iv))
	var block cipher.Block
	var err error
	if block, err = aes.NewCipher([]byte(secretKey)); err != nil {
		log.Println(err)
		return ""
	}
	encrypt := cipher.NewCBCEncrypter(block, iv)
	var source = PKCS5Padding([]byte(origin_data), 16)
	dst := make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)

	return hex.EncodeToString(dst)
}

func AESCBCDecrypt(secretKey SecretKey, encrypt_data string) string {
	iv := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}

	var block cipher.Block
	var err error
	if block, err = aes.NewCipher([]byte(secretKey)); err != nil {
		log.Println(err)
		return ""
	}
	encrypt := cipher.NewCBCDecrypter(block, iv)

	var source []byte

	if source, err = hex.DecodeString(encrypt_data); err != nil {
		log.Println(err)
		return ""
	}
	dst := make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)

	return string(PKCS5Unpadding(dst))
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length-unpadding <= 0 {
		return []byte{}
	}
	return origData[:(length - unpadding)]
}
