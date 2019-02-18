package tool

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type Hashids struct {
}

func (Hashids) EncodeShareKey(UserID uint64) string {
	return "UserID:" + strconv.Itoa(int(UserID))
}
func (Hashids) DecodeShareKey(ShareKey string) uint64 {
	_ShareKey, _ := url.QueryUnescape(ShareKey)
	SuperiorID, _ := strconv.ParseUint(strings.Split(_ShareKey, ":")[1], 10, 64)
	return SuperiorID
}

/*func (Hashids)Encode(id uint64) string  {
	tem:=make([]byte,8)
	bytesBuffer := bytes.NewBuffer(tem)
	binary.Write(bytesBuffer, binary.BigEndian, id)
	bb:=bytesBuffer.Bytes()
	mathrand.Seed(time.Now().UnixNano())
	for i:=0;i<8;i++{
		bb[i]=byte(mathrand.Int31n(int32(256)-int32(bb[i+8])))
		bb[i+8] = bb[i+8]+bb[i]
	}
	cc:=hex.EncodeToString(bb)
	return cc
}
func (Hashids)Decode(id string) uint64  {
	bb,err:=hex.DecodeString(id)
	if err!=nil{
		return 0
	}
	for i:=0;i<8;i++{
		pw:=bb[i+8]-bb[i]
		if pw>255{
			pw = 0
		}
		bb[i+8] = pw
		bb[i]=0
	}
	bytesBuffer:= bytes.NewBuffer(bb[8:])
	var ii uint64
	binary.Read(bytesBuffer, binary.BigEndian, &ii)
	return ii
}*/

//const public_PassWord = "96E5F29353C4A335D2FC4A71DFC8DA3D" // 公共加密字符串
const public_PassWord = "96E5F29353C4A335D2FC4A71DFC8DA3D" // 公共加密字符串

//加密
func CipherEncrypter(tkey, tvalue string) string {
	key := []byte(tkey)
	plaintext := []byte(tvalue)

	BlockSize := aes.BlockSize

	block, err := aes.NewCipher(key)
	if err != nil {
		CheckError(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, BlockSize+len(plaintext))
	iv := ciphertext[:BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		CheckError(err)
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
func CipherDecrypterData(source string) string {

	str := CipherDecrypter(public_PassWord, source)
	return str
}
func CipherEncrypterData(source string) string {
	str := CipherEncrypter(public_PassWord, source)
	return str
}
func HMACSha1(text, key string) []byte {
	keyByte := []byte(key)
	mac := hmac.New(sha1.New, keyByte)
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


