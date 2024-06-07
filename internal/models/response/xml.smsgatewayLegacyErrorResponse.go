package response

import (
	"encoding/xml"
)

type MessageErrorLegacyResponse struct {
	XMLName xml.Name `xml:"env:Envelope"`
	Env     string   `xml:"xmlns:env,attr"`
	Header  struct{} `xml:"env:Header"`
	Body    FaultLegacyBody
}

// Header representa un header vacío en SOAP
type Header struct{}

// FaultBody representa el cuerpo del mensaje SOAP que incluye la estructura de Fault
type FaultLegacyBody struct {
	XMLName xml.Name `xml:"env:Body"`
	Fault   FaultLegacy
}

// Fault define la estructura del mensaje de error SOAP
type FaultLegacy struct {
	XMLName     xml.Name `xml:"env:Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
}

// GenerateFaultXML crea y serializa un XML de tipo Fault
func GenerateFaultXML(code string, message string) string {
	fault := FaultLegacy{
		FaultCode:   code,
		FaultString: message,
	}
	body := FaultLegacyBody{Fault: fault}
	env := MessageErrorLegacyResponse{
		Env:    "http://schemas.xmlsoap.org/soap/envelope/",
		Header: Header{},
		Body:   body,
	}

	output, _ := xml.MarshalIndent(env, "", "    ")

	return string(output)
}
