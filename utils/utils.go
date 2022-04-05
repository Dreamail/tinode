package utils

import "encoding/base64"

func Base64encode(str string) (result []byte) {
	return []byte(base64.StdEncoding.EncodeToString([]byte(str)))
}
