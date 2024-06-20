package config

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/scch94/Gconfiguration"
	"github.com/scch94/ins_log"
)

var Config MicropagosConfiguration

func Upconfig(ctx context.Context) error {

	//traemos el contexto y le setiamos el contexto actual
	// Agregamos el valor "packageName" al contexto
	ctx = ins_log.SetPackageNameInContext(ctx, "config")

	ins_log.Info(ctx, "starting to get the config struct ")
	err := Gconfiguration.GetConfig(&Config, "../config", "micropagosMsgGatewayConfig.json")

	if err != nil {
		ins_log.Fatalf(ctx, "error in Gconfiguration.GetConfig() ", err)
		return err
	}
	return nil
}

type MicropagosConfiguration struct {
	LogLevel                 string         `json:"log_level"`
	Log_name                 string         `json:"log_name"`
	Client                   Client         `json:"client"`
	GetMask                  EndpointConfig `json:"getMask"`
	GetUserDomain            EndpointConfig `json:"getUserDomain"`
	UpdateUsersLastLogin     EndpointConfig `json:"updateUsersLastLogin"`
	ServPort                 string         `json:"server_port"`
	GetUsersInfo             EndpointConfig `json:"getUsersInfo"`
	GetFilterDatabase        EndpointConfig `json:"getFilterDatabase"`
	InsertMessage            EndpointConfig `json:"inserMessage"`
	SMSGateway               EndpointConfig `json:"smsGateway"`
	Portabilidad             EndpointConfig `json:"portabilidad"`
	MaxMessageLength         int            `json:"max_message_length"`
	MobileRegex              string         `json:"mobil_regex"`
	Raven                    []RavenService `json:"raven"`
	UpdatesTimeInMinutes     int            `json:"updates_time_in_minutes"`
	UpdateLastLoginInMinutes int            `json:"update_last_login_time_in_minutes"`
	UseHarcodeShortNumber    bool           `json:"use_harcode_short_number"`
}
type Client struct {
	MaxIdleConns           int  `json:"maxIdleConns"`
	MaxConnsPerHost        int  `json:"maxConnsPerHost"`
	MaxIdleConnsPerHost    int  `json:"maxIdleConnsPerHost"`
	IdleConnTimeoutSeconds int  `json:"idleConnTimeoutSeconds"`
	DisableCompression     bool `json:"disableCompression"`
	PetitionsTimeOut       int  `json:"petitionsTimeOut"`
	DisableKeepAlives      bool `json:"disableKeepAlives"`
}

type EndpointConfig struct {
	URL      string `json:"url"`
	Method   string `json:"method"`
	Timeout  int    `json:"timeout"`
	Username string `json:"portabilidad_user,omitempty"`
	Password string `json:"portabilidad_password,omitempty"`
}

type RavenService struct {
	Name        string `json:"name"`
	SendMail    bool   `json:"sendMail"`
	ShortNumber int    `json:"shortNumber"`
}

func (m MicropagosConfiguration) ConfigurationString() string {
	configJSON, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("Error al convertir la configuración a JSON: %v", err)
	}
	return string(configJSON)
}

func (m MicropagosConfiguration) SearchShortNumber(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, raven := range m.Raven {
		if strings.ToLower(strings.TrimSpace(raven.Name)) == name {
			return fmt.Sprint(raven.ShortNumber)

		}
	}
	return ""
}
