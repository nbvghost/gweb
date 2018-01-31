package tool

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math"
	rand2 "math/rand"
	"strconv"
	"time"
)

func UUID() string {
	rander := rand.Reader
	uuid := make([]byte, 64)
	if _, err := io.ReadFull(rander, uuid); err != nil {
		panic(err.Error()) // rand should never fail
	}

	t := base64.URLEncoding.EncodeToString(uuid)

	t = t + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatInt(rand2.Int63n(math.MaxInt64), 10)
	//fmt.Println(t)

	h := md5.New()
	h.Write([]byte(t))
	return hex.EncodeToString(h.Sum(nil))
	//return uuid
}
