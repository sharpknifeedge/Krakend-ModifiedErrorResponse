package configs

import (
	"io/ioutil"

	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/utils/pkg/env"
	"gopkg.in/yaml.v2"
)

var globalConfig types.Config

func InitAppConfig() error {

	data, err := ioutil.ReadFile(env.Str("CONFIG_PATH", "./cmd/configs") + "/app.yaml")
	if err != nil {
		return err
	}

	var config types.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	globalConfig = config
	return nil
}

func Get() types.Config {

	return globalConfig
}
