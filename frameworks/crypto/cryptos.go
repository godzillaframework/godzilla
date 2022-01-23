package crypto

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func getRandString(length int) string {
	var container string
	var str = "0123456789abcdefghijklmnopqrstuvwxyz"

	b := bytes.NewBufferString(str)
	len := b.Len()
	bigInt := big.NewInt(int64(len))
	for i := 0; i < length; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}

	return container
}
