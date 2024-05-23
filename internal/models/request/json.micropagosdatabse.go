package request

import (
	"encoding/hex"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
)

type InsertMessageRequest struct {
	Id                            *uint64    `json:"id,omitempty"`
	Type                          *string    `json:"type,omitempty"` // Nombre del campo en MySQL
	Content                       *string    `json:"content,omitempty"`
	MobileNumber                  *string    `json:"mobile_number,omitempty"`
	MobileCountryISOCode          *string    `json:"mobile_country_iso_code,omitempty"`
	ShortNumber                   *string    `json:"short_number,omitempty"`
	Telco                         *string    `json:"telco,omitempty"`
	Created                       *time.Time `json:"created,omitempty"`
	RoutingType                   *string    `json:"routing_type,omitempty"`
	MatchedPattern                *string    `json:"matched_pattern,omitempty"`
	ServiceID                     *string    `json:"service_id,omitempty"`
	TelcoID                       *string    `json:"telco_id,omitempty"`
	SessionAction                 *string    `json:"session_action,omitempty"`
	SessionParametersMap          *string    `json:"session_parameters_map,omitempty"`
	SessionTimeoutSeconds         *uint64    `json:"session_timeout_seconds,omitempty"`
	Priority                      *uint64    `json:"priority,omitempty"`
	ClientID                      *string    `json:"client_id,omitempty"`
	URL                           *string    `json:"url,omitempty"`
	AccessTimeoutSeconds          *uint64    `json:"access_timeout_seconds,omitempty"`
	RequestID                     *uint64    `json:"request_id,omitempty"`
	DefaultActionID               *uint64    `json:"default_action_id,omitempty"`
	ApplicationID                 *uint64    `json:"application_id,omitempty"`
	SessionID                     *uint64    `json:"session_id,omitempty"`
	Processed                     *time.Time `json:"processed,omitempty"`
	MillisSinceRequest            *uint64    `json:"millis_since_request,omitempty"`
	SessionApplicationName        *string    `json:"session_application_name,omitempty"`
	Sendafter                     *string    `json:"sendafter,omitempty"`
	Sendbefore                    *string    `json:"sendbefore,omitempty"`
	Sent                          *time.Time `json:"sent,omitempty"`
	Status                        *string    `json:"status,omitempty"`
	AccessTimeoutHandlerQueuename *string    `json:"access_timeout_handler_queuename,omitempty"`
	UseUnsupportedMobilesRegistry *uint64    `json:"use_unsupported_mobiles_registry,omitempty"`
	OriginName                    *string    `json:"origin_name,omitempty"`
}

func hexaToString(hexa string) string {
	decodedString, err := hex.DecodeString(hexa)
	if err != nil {
		return hexa
	}
	// Convertir la cadena decodificada en una cadena normal
	normalString := string(decodedString)
	return normalString
}

func NewInsertMessageRequest(validationStruct *helper.ToValidate, utfi string) *InsertMessageRequest {
	// Crear la estructura InsertMessageRequest
	var mobileNumber string
	if validationStruct.Result == constants.RESULT_SENT {
		mobileNumber = validationStruct.Mobile[len(validationStruct.Mobile)-8:]
	} else {
		mobileNumber = validationStruct.Mobile
	}
	insertMessageRequest := &InsertMessageRequest{
		Created:                timePtr(validationStruct.StartPetition),
		Type:                   stringPtr(constants.INSERT_TYPE),
		Content:                stringPtr(hexaToString(validationStruct.Message)),
		MobileNumber:           stringPtr(mobileNumber),
		MobileCountryISOCode:   stringPtr(constants.INSERT_MOBILEC_OUNTRY_ISO_CODE),
		ShortNumber:            stringPtr(validationStruct.ShortNumber),
		Telco:                  stringPtr(validationStruct.Telco),
		SessionAction:          stringPtr(constants.INSERT_SESSION_ACTION),
		SessionApplicationName: stringPtr(constants.INSERT_SESSION_APPLICATION_NAME),
		Sendafter:              stringPtr(validationStruct.SendAfter),
		Sendbefore:             stringPtr(validationStruct.SendBefore),
		Sent:                   timePtr(validationStruct.SendMessageTime),
		Processed:              timePtr(time.Now()),
		Status:                 &validationStruct.Result,
		OriginName:             stringPtr(validationStruct.UseOriginName),
	}
	return insertMessageRequest
}

// Funciones auxiliares para crear punteros a valores primitivos
func stringPtr(s string) *string {
	return &s
}
func timePtr(t time.Time) *time.Time {
	return &t
}
