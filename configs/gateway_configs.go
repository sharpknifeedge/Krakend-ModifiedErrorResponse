package configs

import (
	"errors"
	"io/ioutil"

	"gitlab.boomerangapp.ir/back/pg/types"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/env"
	"gopkg.in/yaml.v2"
)

var gateConfigsMap map[string]types.PSPConfig

func InitGatewaysConfig() error {

	data, err := ioutil.ReadFile(env.Str("CONFIG_PATH", "./cmd/configs") + "/gateways_config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &gateConfigsMap)
	if err != nil {
		return err
	}

	return nil
}

func GetGateConfig(gate string) (*types.PSPConfig, bool) {

	val, ok := gateConfigsMap[gate]
	if !ok {
		return nil, false
	}

	return &val, true
}

func GetAllGates() ([]types.PSPResponse, error) {

	var resp []types.PSPResponse
	pspNames := [][]string{
		{consts.PSPSadad, consts.PSPSadadPersion},
		{consts.PSPSep, consts.PSPSepPersion},
		{consts.PSPAtipay, consts.PSPAtipayPeraion},
		{consts.PSPPec, consts.PSPPecPeraion},
	}

	for _, pspName := range pspNames {
		psp, ok := gateConfigsMap[pspName[0]]
		if !ok {
			continue
		}
		pspResp := types.PSPResponse{
			Name:  pspName[0],
			Label: pspName[1],
			Logo:  psp.Logo,
		}
		for _, terminal := range psp.Terminals {
			pspResp.Terminals = append(pspResp.Terminals, uint16(terminal.ID))
		}
		resp = append(resp, pspResp)
	}

	if len(resp) == 0 {
		return nil, errors.New("no psp provider is valid")
	}
	return resp, nil
}
