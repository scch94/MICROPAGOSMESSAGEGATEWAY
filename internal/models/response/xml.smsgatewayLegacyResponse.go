package response

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Envelope define la estructura del documento XML SOAP
type MessageLegacyResponse struct {
	XMLName xml.Name `xml:"env:Envelope"`
	Env     string   `xml:"xmlns:env,attr"`
	Header  struct{} `xml:"env:Header"`
	Body    LegacyBody
}

// Body representa el cuerpo del mensaje SOAP que puede contener diferentes tipos de respuestas
type LegacyBody struct {
	XMLName                     xml.Name                     `xml:"env:Body"`
	SendMassiveMessagesResponse *SendMassiveMessagesResponse `xml:"ns2:sendMassiveMessagesResponse,omitempty"`
	SendResponse                *SendResponse                `xml:"ns2:sendResponse,omitempty"`
}

// SendMassiveMessagesResponse representa una respuesta vacía para sendMassiveMessages
type SendMassiveMessagesResponse struct {
	XMLNS string `xml:"xmlns:ns2,attr"`
}

// SendResponse representa una respuesta vacía para sendResponse
type SendResponse struct {
	XMLNS string `xml:"xmlns:ns2,attr"`
}

// GenerateXML genera un XML basado en el tipo de respuesta especificada
func GenerateLegacyXML(useMassive bool) (string, error) {
	messageLegacyResponse := MessageLegacyResponse{
		Env:    "http://schemas.xmlsoap.org/soap/envelope/",
		Header: Header{},
	}

	if useMassive {
		messageLegacyResponse.Body = LegacyBody{
			SendMassiveMessagesResponse: &SendMassiveMessagesResponse{
				XMLNS: "http://webservices.ravenws.micropagos.com.uy/",
			},
		}
	} else {
		messageLegacyResponse.Body = LegacyBody{
			SendResponse: &SendResponse{
				XMLNS: "http://webservices.ravenws.micropagos.com.uy/",
			},
		}
	}

	output, err := xml.MarshalIndent(messageLegacyResponse, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error al generar XML: %s", err)
	}

	// Reemplaza <env:Header></env:Header> por <env:Header/>
	outputStr := strings.Replace(string(output), "<env:Header></env:Header>", "<env:Header/>", 1)
	if useMassive {
		outputStr = strings.Replace(string(outputStr), `<ns2:sendMassiveMessagesResponse xmlns:ns2="http://webservices.ravenws.micropagos.com.uy/"></ns2:sendMassiveMessagesResponse>`, `<ns2:sendMassiveMessagesResponse xmlns:ns2="http://webservices.ravenws.micropagos.com.uy/"/>`, 1)
	} else {
		outputStr = strings.Replace(string(outputStr), `<ns2:sendResponse xmlns:ns2="http://webservices.ravenws.micropagos.com.uy/"></ns2:sendResponse>`, `<ns2:sendResponse xmlns:ns2="http://webservices.ravenws.micropagos.com.uy/"/>`, 1)
	}
	return xml.Header + outputStr, nil
}
