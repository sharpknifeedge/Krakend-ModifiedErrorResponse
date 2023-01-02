package validation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var boilerPkg = "entities"

//SetSchemaPkg when using sqlBoiler.
//It's already set "entities"
func SetSchemaPkg(pkg string) {

	boilerPkg = pkg
}

type MysqlValidationErr struct {
	field, message string
}

func (e *MysqlValidationErr) Error() string {

	if e.field != "" {
		return fmt.Sprintf(e.message, e.field)
	}

	return e.message
}

//ValidateByErr can validate sqlBoiler and Mysql errors
func ValidateByErr(err error) error {

	if err.Error() == "sql: no rows in result set" {
		return &MysqlValidationErr{
			message: "داده ای یافت نشد",
		}
	}

	var errorMessage, errorNumber string
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		errorNumber = strconv.Itoa(int(mysqlErr.Number))
		mysqlErr.Message = errorMessage

	} else if msg, ok := isBoilerError(err); ok {
		strs := strings.Split(msg, ":")
		nn := strings.Split(strs[2], " ")
		errorNumber = nn[len(nn)-1]
		errorMessage = strs[len(strs)-1]

	} else {
		return nil
	}

	switch errorNumber {
	case "1062":
		words := strings.Split(errorMessage, " ")
		lastWord := strings.TrimRight(words[len(words)-1], "'")
		return &MysqlValidationErr{
			field:   lastWord[strings.IndexByte(lastWord, '.')+1:],
			message: "فیلد %s تکراری است",
		}

	case "1406":
		return &MysqlValidationErr{
			field: errorMessage[strings.
				IndexByte(errorMessage, '\'')+1 : strings.
				LastIndexByte(errorMessage, '\'')],
			message: "فیلد %s طولانی است",
		}

	case "1048":
		return &MysqlValidationErr{
			field: errorMessage[strings.
				IndexByte(errorMessage, '\'')+1 : strings.
				LastIndexByte(errorMessage, '\'')],
			message: "فیلد %s نمیتواند خالی باشد",
		}

	case "1265":
		return &MysqlValidationErr{
			field: errorMessage[strings.
				IndexByte(errorMessage, '\'')+1 : strings.
				LastIndexByte(errorMessage, '\'')],
			message: "مقدار غیر مجاز برای فیلد %s",
		}

	case "1216":
		return &MysqlValidationErr{
			message: "داده ناهماهنگ!",
		}
	}

	return nil
}

func isBoilerError(err error) (string, bool) {

	errMsg := err.Error()
	firstElem := strings.Split(errMsg, ":")[0]
	if firstElem != boilerPkg && len(errMsg) > len(firstElem)+2 {
		errMsg = errMsg[len(firstElem)+2:]
	}

	return errMsg, strings.Split(errMsg, ":")[0] == boilerPkg
}
