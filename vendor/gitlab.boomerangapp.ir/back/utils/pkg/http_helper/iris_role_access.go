package http_helper

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"gitlab.boomerangapp.ir/back/utils/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/response"
)

var unauthorizedResp = response.ErrorResponse(response.Unauthorized)

func writeUnauthorizeResp(ctx iris.Context) {
	ctx.StatusCode(iris.StatusUnauthorized)
	ctx.JSON(&unauthorizedResp)
}

func Authorize(forceUserID bool, roles ...string) func(iris.Context) {
	return func(ctx iris.Context) {
		userIDString := ctx.GetHeader(consts.HTTPUserID)
		if len(userIDString) > 0 {
			uid, _ := strconv.ParseUint(userIDString, 10, 64)
			ctx.Values().Set(consts.HTTPUserID, uid)
		} else if forceUserID {
			writeUnauthorizeResp(ctx)
			return
		}

		if len(roles) > 0 {
			r := ctx.GetHeader(consts.HTTPROleHeader)
			seen := false
			for _, role := range roles {
				if r == role {
					seen = true
					break
				}
			}
			if !seen {
				writeUnauthorizeResp(ctx)
			}
		}

		ctx.Next()
	}
}
