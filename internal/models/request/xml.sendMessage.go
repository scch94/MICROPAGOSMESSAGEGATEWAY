package request

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
)

type SendMessageRequest struct {
	XMLName xml.Name
	Body    Body
}

type Body struct {
	SendMassiveMessages SendMassiveMessages `xml:"sendMassiveMessages,omitempty"`
	Send                Send                `xml:"send,omitempty"`
}

type SendMassiveMessages struct {
	MassiveMessage   string             `xml:"massiveMessage,omitempty"`   // Campo opcional
	MobileMessageDto []MobileMessageDto `xml:"mobileMessageDto,omitempty"` // Campo opcional
	SendAfter        string             `xml:"sendAfter,omitempty"`        // Campo opcional
	SendBefore       string             `xml:"sendBefore,omitempty"`       // Campo opcional
	ShortNumber      string             `xml:"shortNumber,omitempty"`      // Campo opcional
	UseOriginName    string             `xml:"useOriginName,omitempty"`    // Campo opcional
	Priority         string             `xml:"priority,omitempty"`         // Campo opcional
}

type MobileMessageDto struct {
	Message string `xml:"message,omitempty"` // Campo opcional
	Mobile  string `xml:"mobile,omitempty"`  // Campo opcional
}

type Send struct {
	Mobile        string `xml:"mobile,omitempty"`        // Campo opcional
	Message       string `xml:"message,omitempty"`       // Campo opcional
	UseOriginName string `xml:"useOriginName,omitempty"` // Campo opcional
}

// Método común para imprimir la estructura de datos
func (e *SendMessageRequest) DetermineRequestType() (string, error) {
	if e.Body.SendMassiveMessages.MobileMessageDto != nil {
		return constants.MASSIVE_MESSAGE, nil

	} else if e.Body.Send.Mobile != "" {
		return constants.SEND, nil
	} else {
		return "", errors.New("peticion no valida")
	}
}
func (s *Send) SendToString() string {
	return fmt.Sprintf("Mobile: %s, Message: %s, UseOriginName: %s", s.Mobile, s.Message, s.UseOriginName)
}
func (s *SendMassiveMessages) SendMassiveMessagesToString() string {
	result := fmt.Sprintf("MassiveMessage: %s, SendAfter: %s, SendBefore: %s, ShortNumber: %s, UseOriginName: %s, Priority: %s\n",
		s.MassiveMessage, s.SendAfter, s.SendBefore, s.ShortNumber, s.UseOriginName, s.Priority)

	for i, dto := range s.MobileMessageDto {
		result += fmt.Sprintf("MobileMessageDto[%d]: {Message: %s, Mobile: %s}\n", i, dto.Message, dto.Mobile)
	}

	return result
}

func (s *SendMessageRequest) SendToValidation(r *SmsGatewayRequest) {

}
