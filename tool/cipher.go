package tool

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

const public_PassWord = "96E5F29353C4A335D2FC4A71DFC8DA3D" // 公共加密字符串

//加密
func CipherEncrypter(tkey, tvalue string) string {
	key := []byte(tkey)
	plaintext := []byte(tvalue)

	block, err := aes.NewCipher(key)
	if err != nil {
		CheckError(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		CheckError(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	//fmt.Printf("%s", ciphertext)
	return hex.EncodeToString(ciphertext)

}
func CipherDecrypterData(source string) string {

	str := CipherDecrypter(public_PassWord, source)
	return str
}
func CipherEncrypterData(source string) string {
	str := CipherEncrypter(public_PassWord, source)
	return str
}
func Md5(valeu string) string {
	ddf := md5.New()
	ddf.Write([]byte(valeu))
	md5Str := hex.EncodeToString(ddf.Sum(nil))
	return strings.ToUpper(md5Str)
}

//解密
func CipherDecrypter(tkey string, crypter string) string {
	key := []byte(tkey)
	ciphertext, _ := hex.DecodeString(crypter)

	block, err := aes.NewCipher(key)
	if err != nil {
		CheckError(err)
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		CheckError(errors.New("必须是aes.BlockSize的倍数"))
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	//fmt.Printf("%s", ciphertext)
	return string(ciphertext)
}
