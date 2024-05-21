package response

import (
	"encoding/xml"
	"fmt"
)

type MessageResponse struct {
	XMLName xml.Name `xml:"env:Envelope"`
	Env     string   `xml:"xmlns:env,attr"`
	Header  string   `xml:"env:Header"`
	Body    Body
}

// Body representa el cuerpo del mensaje SOAP
type Body struct {
	XMLName xml.Name `xml:"env:Body"`
	Success []Success
	Error   []Fault
}

// Fault representa la estructura de un mensaje de error SOAP
type Success struct {
	XMLName    xml.Name `xml:"env:Success"`
	DatabaseID string   `xml:"database_id"`
	Code       string   `xml:"code"`
	Message    string   `xml:"message"`
}
type Fault struct {
	XMLName    xml.Name `xml:"env:Fault"`
	DatabaseID string   `xml:"database_id"`
	Status     string   `xml:"status"`
	Message    string   `xml:"message"`
}

type Result struct {
	Code       string
	Message    string
	DatabaseID string
}

// NewMessageResponse crea, serializa y devuelve una estructura MessageResponse con el FaultCode y FaultString especificados
func GenerateXML(results []Result) string {
	body := Body{}

	for _, result := range results {
		if result.Code == "0" {
			body.Success = append(body.Success, Success{
				Code:       result.Code,
				Message:    result.Message,
				DatabaseID: result.DatabaseID,
			})
		} else {
			body.Error = append(body.Error, Fault{
				Status:     result.Code,
				Message:    result.Message,
				DatabaseID: result.DatabaseID,
			})
		}
	}

	env := MessageResponse{
		Env:    "http://schemas.xmlsoap.org/soap/envelope/",
		Header: "",
		Body:   body,
	}

	output, err := xml.MarshalIndent(env, "", "    ")
	if err != nil {
		return fmt.Sprintf("Error al generar XML: %s", err)
	}

	return string(output)
}
