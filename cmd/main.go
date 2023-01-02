package main

import (
	"gitlab.boomerangapp.ir/back/pg/configs"
	globalInit "gitlab.boomerangapp.ir/back/pg/init"
	"gitlab.boomerangapp.ir/back/pg/internal/version"
	internalHTTP "gitlab.boomerangapp.ir/back/pg/web/http"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

func main() {

	//NOTE: for development
	// proxyUrl, _ := url.Parse("socks5://:8989")
	// http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	//init sentry
	sentrylog.Init("", version.Version)
	defer sentrylog.Recover()

	//initialization
	err := globalInit.Init()
	if err != nil {
		sentrylog.Fatal(err)
	}

	//run http api
	sentrylog.Info("HTTP api started at port " + configs.Get().AppConfigs.HttpPort)
	err = internalHTTP.RunHttpApi()
	if err != nil {
		sentrylog.Fatal(err)
	}
}
