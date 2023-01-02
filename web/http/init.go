package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/pg/internal/version"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/http_helper"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

var httpServer http.Server

func InitHttpListener() {

	router := mux.NewRouter()
	http_helper.MuxHelthCheck(router, version.Version)
	router = router.PathPrefix(consts.HttpPrefix).Subrouter()
	router.NotFoundHandler = http.HandlerFunc(http_helper.MuxNotFound)
	router.Use(sentrylog.SentryMiddleware)

	router.Methods("POST").
		Path("/request").
		HandlerFunc(handleRequest)
	router.Methods("POST").
		Path(consts.ReturnURL + "/{paymentID}").
		HandlerFunc(handleConfirm)
	router.Methods("GET").
		Path(consts.RouteRedirect + "/{psp}/{token}").
		HandlerFunc(handleRedirect)
	router.Methods("GET").
		Path("/psp").
		HandlerFunc(handleGetPSPs)
	router.Methods("GET").
		Path("/transactions").
		HandlerFunc(handleTransactions)

	timeout := configs.Get().AppConfigs.HttpReadTimeout
	httpServer = http.Server{
		ReadTimeout: time.Duration(timeout) * time.Second,
		Addr:        configs.Get().AppConfigs.HttpPort,
		Handler:     router,
	}
}

func RunHttpApi() error {

	err := httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
