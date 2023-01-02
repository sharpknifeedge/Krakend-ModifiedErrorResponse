package init

import (
	"gitlab.boomerangapp.ir/back/pg/configs"
	"gitlab.boomerangapp.ir/back/pg/internal/db"
	"gitlab.boomerangapp.ir/back/pg/web/http"
)

func Init() error {

	//init global configs
	err := configs.InitAppConfig()
	if err != nil {
		return err
	}

	//init gateway configs
	err = configs.InitGatewaysConfig()
	if err != nil {
		return err
	}

	//init DB
	err = db.InitDB()
	if err != nil {
		return err
	}

	//init http api
	http.InitHttpListener()

	return nil
}
