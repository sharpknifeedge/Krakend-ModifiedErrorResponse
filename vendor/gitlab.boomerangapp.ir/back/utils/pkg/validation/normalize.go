package validation

import (
	"strconv"
	"strings"
)

var persianDigits = []string{"۰", "۱", "۲", "۳", "۴", "۵", "۶", "۷", "۸", "۹"}

//PersianToEnglishDigits ,Convert Persian numbers to English
func PersianToEnglishDigits(text string) string {
	if len(text) == 0 {
		return ""
	}
	for k, i := range persianDigits {
		text = strings.ReplaceAll(text, i, strconv.Itoa(k))
	}
	return text
}

//EnglishToPersianDigits ,Convert English numbers to Persian
func EnglishToPersianDigits(text string) string {
	if len(text) == 0 {
		return ""
	}
	for k, i := range persianDigits {
		text = strings.ReplaceAll(text, strconv.Itoa(k), i)
	}
	return text
}
