package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/nbvghost/glog"
	"io"
	"strings"
)

type SecretKey string

func NewSecretKey(key string) SecretKey {
	return SecretKey(Md5ByString(key))
}

//加密
func CipherEncrypter(tkey SecretKey, tvalue string) string {
	key := []byte(tkey)
	plaintext := []byte(tvalue)

	BlockSize := aes.BlockSize

	block, err := aes.NewCipher(key)
	if err != nil {
		glog.Error(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, BlockSize+len(plaintext))
	iv := ciphertext[:BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		glog.Error(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	//fmt.Printf("%s", ciphertext)
	return hex.EncodeToString(ciphertext)

}

//解密
func CipherDecrypter(tkey SecretKey, crypter string) string {
	key := []byte(tkey)
	ciphertext, _ := hex.DecodeString(crypter)

	block, err := aes.NewCipher(key)
	if err != nil {
		glog.Error(err)
		return ""
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		glog.Error(errors.New("必须是aes.BlockSize的倍数"))
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

func HMACSha1(text, key string) []byte {
	keyByte := []byte(key)
	mac := hmac.New(sha1.New, keyByte)
	mac.Write([]byte(text))
	return mac.Sum(nil)
}
func HMACSha256(text, key string) []byte {
	keyByte := []byte(key)
	mac := hmac.New(sha256.New, keyByte)
	mac.Write([]byte(text))
	return mac.Sum(nil)
}
func Sha1ByBytes(value []byte) string {
	sha1Sign := sha1.New()
	//fmt.Println(list.Join("|") + "|cd99858d693d38104e7df5c7f771f474")
	sha1Sign.Write(value)
	mySign := hex.EncodeToString(sha1Sign.Sum(nil))
	return strings.ToUpper(mySign)
}

func Sha1ByString(value string) string {
	mySign := Sha1ByBytes([]byte(value))
	return mySign
}
func Md5ByString(valeu string) string {
	ddf := md5.New()
	ddf.Write([]byte(valeu))
	md5Str := hex.EncodeToString(ddf.Sum(nil))
	return strings.ToUpper(md5Str)
}
func Md5By16String(valeu string) string {
	ddf := Md5ByString(valeu)
	return strings.ToUpper(string(ddf[8:24]))
}
func Md5ByBytes(valeu []byte) string {
	ddf := md5.New()
	ddf.Write(valeu)
	md5Str := hex.EncodeToString(ddf.Sum(nil))
	return strings.ToUpper(md5Str)
}
