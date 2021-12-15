package toolib

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Hash(bys []byte) string {
	hash := md5.New()
	hash.Write(bys)
	return hex.EncodeToString(hash.Sum(nil))
}
