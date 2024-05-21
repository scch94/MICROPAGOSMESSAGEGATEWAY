package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scch94/Gconfiguration"
	"github.com/scch94/ins_log"
)

//lint:ignore SA1029 "Using built-in type string as key for context value intentionally"
var ctx = context.WithValue(context.Background(), "packageName", "config")
var Config MicropagosConfiguration

func Upconfig() error {
	ins_log.Info(ctx, "starting to get the config struct ")
	err := Gconfiguration.GetConfig(&Config)
	if err != nil {
		ins_log.Fatalf(ctx, "error in Gconfiguration.GetConfig() ", err)
		return err
	}
	return nil
}

type MicropagosConfiguration struct {
	GetUserDomain     EndpointConfig `json:"getUserDomain"`
	ServPort          string         `json:"server_port"`
	GetFilterDatabase EndpointConfig `json:"getFilterDatabase"`
	InsertMessage     EndpointConfig `json:"inserMessage"`
	SMSGateway        EndpointConfig `json:"smsGateway"`
	Portabilidad      EndpointConfig `json:"portabilidad"`
	MaxMessageLength  int            `json:"max_message_length"`
	MobileRegex       string         `json:"mobil_regex"`
	Raven             []RavenService `json:"raven"`
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
	for _, raven := range m.Raven {
		if raven.Name == name {
			return fmt.Sprint(raven.ShortNumber)

		}
	}
	return ""
}