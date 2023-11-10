package helpers

import (
	"check-price/src/common/configs"
	"encoding/base64"
)

func EncryptText(plainText string) string {
	appKey := configs.Get().AppKey
	keyLength := len(appKey)
	plainTextLength := len(plainText)
	maxLength := max(plainTextLength, keyLength)
	var result []byte

	for i := 0; i < maxLength; i++ {
		if i < plainTextLength {
			result = append(result, plainText[i])
		}
		if i < keyLength {
			result = append(result, appKey[i])
		}
	}
	return base64.StdEncoding.EncodeToString(result)
}

func DecryptText(mixedString string) string {
	decoded, _ := base64.StdEncoding.DecodeString(mixedString)
	var str1 []byte
	var str2 []byte

	for i := 0; i < len(decoded); i += 2 {
		str1 = append(str1, decoded[i])
		if i+1 < len(decoded) {
			str2 = append(str2, decoded[i+1])
		}
	}
	return string(str1)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
