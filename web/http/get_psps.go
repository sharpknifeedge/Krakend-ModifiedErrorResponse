package http

import (
	"net/http"

	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/utils/pkg/http_helper"
	"gitlab.boomerangapp.ir/back/utils/pkg/response"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

func handleGetPSPs(w http.ResponseWriter, r *http.Request) {

	//init response handler
	mr := http_helper.MuxResp{
		Writer: w,
		Hub:    sentrylog.GetHubFromCtx(r.Context()),
	}

	//load psp configs
	gateways, err := configs.GetAllGates()
	if err != nil {
		mr.HandleErr(err, response.InternalError)
		return
	}

	//send response
	mr.SendResp(gateways, nil, false)
}
