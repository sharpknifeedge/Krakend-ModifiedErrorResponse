package types

type (
	Config struct {
		AppConfigs AppConfigs `yaml:"app_configs"`
		DBConfigs  DBConfigs  `yaml:"db_configs"`
	}
	AppConfigs struct {
		HttpPort        string `yaml:"http_port"`
		BaseURL         string `yaml:"base_url"`
		GrpcUserAddr    string `yaml:"grpc_user_addr"`
		HtmlSuccess     string `yaml:"html_success"`
		Htmlfailed      string `yaml:"html_failed"`
		HttpReadTimeout uint8  `yaml:"http_read_timeout"`
	}
	DBConfigs struct {
		Mysql Mysql `yaml:"mysql"`
	}
	Mysql struct {
		Host     string `yaml:"host"`
		DB       string `yaml:"db"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
)

type PSPConfig struct {
	Terminals []struct {
		TerminalID string `yaml:"terminalID"`
		ApiKey     string `yaml:"apiKey"`
		ID         uint8  `yaml:"id"`
	} `yaml:"terminals"`
	MerchantID string `yaml:"merchantID"`
	Logo       string `yaml:"logo"`
}

type PaymentRequest struct {
	Callback     string `json:"callback"`
	PSP          string `json:"psp"`
	WalletID     uint64 `json:"wallet_id"`
	ID           uint
	Amount       uint  `json:"amount"`
	GatewayID    uint8 `json:"gateway_id"`
	AutoRedirect bool  `json:"auto_redirect"`
}

type PSPResponse struct {
	Name      string   `json:"name"`
	Label     string   `json:"label"`
	Logo      string   `json:"logo"`
	Terminals []uint16 `json:"terminals"`
}

type GateInterface interface {
	GetToken(pr PaymentRequest, mobile string) (string, error)
	Verify(verifyCode string, verifyResp interface{}) error
	MakeRedirectURL(pspToken string) string
	SerializeToDB(resp interface{}) string
	CheckConfirmResp(confirmResp map[string]string) (string, bool)
	CheckVerifyResp(verifyResp interface{}) bool
}
