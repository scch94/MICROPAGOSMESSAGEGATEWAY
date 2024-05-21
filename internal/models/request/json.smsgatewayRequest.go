package request

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
)

type SmsGatewayRequest struct {
	Utfi        string `json:"utfi"`
	ServiceType string `json:"serviceType"`
	//OriginTon            uint8  `json:"origin.ton,omitempty"`
	//OriginNpi            uint8  `json:"origin.npi,omitempty"`
	OriginNumber string `json:"origin.number"`
	//DestinationTon       uint8  `json:"destination.ton,omitempty"`
	//DestinationNpi       uint8  `json:"destination.npi,omitempty"`
	DestinationNumber    string `json:"destination.number"`
	ValidityPeriod       string `json:"validity_period"`
	ScheduleDeliveryTime string `json:"schedule_delivery_time"`
	ProtocolID           uint8  `json:"protocol_id"`
	EsmeClass            uint8  `json:"esmeClass"`
	PriorityFlag         uint8  `json:"priority_flag"`
	RegisteredDelivery   uint8  `json:"registered_delivery"`
	ReplaceIfPresentFlag uint8  `json:"replace_if_present_flag"`
	Data                 string `json:"data"`
	DataHeaderIndicator  uint8  `json:"data_header_indicator"`
	DataCodingScheme     uint8  `json:"data_coding_scheme"`
	DataLength           uint16 `json:"data_length"`
	MessageType          uint8  `json:"messagetype"`
	TLVTag               int    `json:"TLV_tag"`
	TLVLength            int    `json:"TLV_length"`
	TLVValue             string `json:"TLV_value"`
}

func NewSmsGatewayRequest(h helper.ToValidate, utfi string) *SmsGatewayRequest {
	smsGatewayRequest := &SmsGatewayRequest{
		Utfi:                 utfi,
		ServiceType:          "",
		OriginNumber:         h.ShortNumber,
		DestinationNumber:    h.Mobile,
		ValidityPeriod:       "",
		ScheduleDeliveryTime: "",
		ProtocolID:           0,
		EsmeClass:            0,
		PriorityFlag:         3,
		RegisteredDelivery:   0,
		ReplaceIfPresentFlag: 0,
		Data:                 strings.ToUpper(hex.EncodeToString([]byte(h.Message))),
		DataHeaderIndicator:  0,
		DataCodingScheme:     0,
		DataLength:           uint16(len(h.Message)),
		MessageType:          4,
		TLVTag:               515,
		TLVLength:            0,
		TLVValue:             "",
	}

	return smsGatewayRequest

}

func (payload *SmsGatewayRequest) ToJSON() (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
