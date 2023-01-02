package http_helper

import (
	"encoding/json"
	"reflect"

	"github.com/kataras/iris/v12"
	"gitlab.boomerangapp.ir/back/utils/pkg/response"
)

const BodyValue = "body-value"

func ReadBodyJSON(instanceOf interface{}) func(iris.Context) {
	v := reflect.ValueOf(instanceOf)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return func(ctx iris.Context) {
		sPtr := reflect.New(v.Type())
		body, err := ctx.GetBody()
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			resp := response.ErrorResponse(response.InternalError)
			ctx.JSON(&resp)
			return
		}

		if err := json.Unmarshal(body, sPtr.Interface()); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			resp := response.ErrorResponse(response.InternalError)
			ctx.JSON(&resp)
			return
		}

		ctx.Values().Set(BodyValue, sPtr.Interface())
		ctx.Next()
	}
}

type Validator interface {
	Validate() error
}

func ValidateBody(ctx iris.Context) {
	body := ctx.Values().Get(BodyValue)
	if body != nil {
		if validator, ok := body.(Validator); ok {
			if ve := validator.Validate(); ve != nil {
				resp := response.ErrorResponse(response.UnprocessableEntity)
				resp.Message.Body = ve.Error()
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(&resp)
				return
			}
		}
	}
	ctx.Next()
}
