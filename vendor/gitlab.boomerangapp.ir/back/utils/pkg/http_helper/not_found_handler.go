package http_helper

import (
	"encoding/json"
	"net/http"

	"github.com/kataras/iris/v12"
	"gitlab.boomerangapp.ir/back/utils/pkg/response"
)

var notFoundResp = response.ErrorResponse(response.NotFound)

func IrisNotFound(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNotFound)
	ctx.JSON(notFoundResp)
}

func MuxNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(notFoundResp)
}
