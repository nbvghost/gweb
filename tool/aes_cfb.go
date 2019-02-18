package tool

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func AESCFBEndEncrypt(origin_data string)string{

	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	//key, _ := hex.DecodeString("6368616e676520746869732070617373")
	key, _ := hex.DecodeString(public_PassWord)
	plaintext := []byte(origin_data)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	/*if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}*/

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	//fmt.Printf("%x\n", ciphertext[aes.BlockSize:])
	return hex.EncodeToString(ciphertext[aes.BlockSize:])
}

func AESCFBEndDecrypt(encrypt_data string)string{
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString(public_PassWord)
	ciphertext, _ := hex.DecodeString(encrypt_data)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.

	//iv := ciphertext[:aes.BlockSize]
	iv := make([]byte,aes.BlockSize)
	//ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	//fmt.Printf("%s", ciphertext)
	return string(ciphertext)
}
