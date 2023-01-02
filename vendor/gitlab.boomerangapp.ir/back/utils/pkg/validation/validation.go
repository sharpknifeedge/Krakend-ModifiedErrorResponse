package validation

import (
	"log"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func CorrectionNationalCode(code string) (string, bool) {
	reg, err := regexp.Compile("/[^0-9]/")
	if err != nil {
		log.Fatalln(err)
	}
	code = PersianToEnglishDigits(code)
	code = reg.ReplaceAllString(code, "")
	if len(code) != 10 {
		return code, false
	}
	codes := strings.Split(code, "")
	last, _ := strconv.Atoi(codes[9])
	i := 10
	sum := 0
	for in, el := range codes {
		temp, err := strconv.Atoi(el)
		if err != nil {
			log.Fatal(err)
		}
		if in == 9 {
			break
		}
		sum += temp * i
		i -= 1
	}
	mod := sum % 11
	if mod >= 2 {
		mod = 11 - mod
	}
	return code, mod == last
}

func CorrectionEmail(email string) (string, bool) {
	email = PersianToEnglishDigits(email)
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'+/=?^_`{|}~-]+" +
		"@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:" +
		"\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)$")
	if !re.Match([]byte(email)) {
		return email, false
	}
	email = strings.ToLower(email)

	return email, true
}

func CorrectionPhoneNumber(phone string) (string, bool) {
	phone = PersianToEnglishDigits(phone)
	re := regexp.MustCompile(`^(0|\\+98)[0-9]{10}$`)
	if !re.Match([]byte(phone)) {
		return phone, false
	}
	if strings.HasPrefix(phone, "+98") {
		phone = strings.Replace(phone, "+98", "0", 1)
	}
	return phone, true
}

func TrimValidation(v interface{}) {
	vo := reflect.ValueOf(v)
	if vo.Kind() != reflect.Ptr {
		return
	}
	vo = vo.Elem()
	for i := 0; i < vo.NumField(); i++ {
		fieldValue := vo.Field(i)
		if fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.String {
			stringValue := fieldValue.Interface().(*string)
			if stringValue == nil {
				continue
			}
			*stringValue = strings.TrimSpace(*stringValue)
			if *stringValue == "" {
				fieldValue.Set(reflect.Zero(fieldValue.Type()))
			} else {
				fieldValue.Set(reflect.ValueOf(stringValue))
			}
		}
	}
}

func BankCardValidation(text string) bool {
	if len(text) != 16 {
		return false
	}
	var cardTotal int64 = 0
	for i, ch := range text {
		c, err := strconv.ParseInt(string(ch), 10, 8)
		if err != nil {
			return false
		}
		if i%2 == 0 {
			if c*2 > 9 {
				cardTotal = cardTotal + (c * 2) - 9
			} else {
				cardTotal = cardTotal + (c * 2)
			}
		} else {
			cardTotal += c
		}
	}
	return cardTotal%10 == 0
}

func ShebaValidation(text string) bool {
	matched, err := regexp.MatchString("^IR[0-9]{24}$", text)
	if err != nil || !matched {
		return false
	}
	firstForChars := ""
	sheba := ""
	for index := range text {
		if index < 4 {
			if string(text[index]) == "I" {
				firstForChars = firstForChars + "18"
			} else if string(text[index]) == "R" {
				firstForChars = firstForChars + "27"
			} else {
				firstForChars = firstForChars + string(text[index])
			}
		} else {
			sheba = sheba + string(text[index])
		}
	}
	sheba = sheba + firstForChars
	bigNum97 := big.NewInt(97)
	n := new(big.Int)
	shebaNumber, ok := n.SetString(sheba, 10)
	if !ok {
		return false
	}
	bigNum1 := new(big.Int)
	bigNum1.SetInt64(1)

	return 0 == bigNum1.Cmp(shebaNumber.Mod(shebaNumber, bigNum97))
}
