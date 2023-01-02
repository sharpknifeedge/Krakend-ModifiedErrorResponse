package http_helper

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kataras/iris/v12"
	"gitlab.boomerangapp.ir/back/utils/pkg/slices"
)

var ErrorInvalidMimeType = errors.New("فرمت فایل نامعتبر است")

func IrisUpload(c iris.Context, inputName string, path string, fileName string, ValidMimeTypes *[]string) (string, error) {
	_, file, err := c.FormFile(inputName)

	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	headerByte := make([]byte, 512)
	_, err = src.ReadAt(headerByte, 0)
	if err != nil {
		return "", err
	}

	//check mime type if forced
	if ValidMimeTypes != nil {
		if slices.StringContains(*ValidMimeTypes, http.DetectContentType(headerByte)) < 0 {
			return "", ErrorInvalidMimeType
		}
	}

	finalFileName := ""
	if len(fileName) == 0 {
		finalFileName = file.Filename
	} else {
		sl := strings.Split(file.Filename, ".")
		finalFileName = fileName + "." + sl[len(sl)-1]
	}

	//Replace ..
	finalFileName = strings.Replace(finalFileName, "..", "", -1)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755) // #nosec
		if err != nil {
			return "", err
		}
	}
	// Destination
	dst, err := os.Create(path + finalFileName)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return finalFileName, nil
}
