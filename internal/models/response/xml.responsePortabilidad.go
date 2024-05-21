package response

import "encoding/xml"

type PortabilidadResponse struct {
	XMLName xml.Name     `xml:"Envelope"`
	Body    BodyResponse `xml:"Body"`
}

type BodyResponse struct {
	GetTelcoResponse GetTelcoResponse `xml:"getTelcoResponse"`
}

type GetTelcoResponse struct {
	Return Return `xml:"return"`
}

type Return struct {
	TelcoCode int    `xml:"telcoCode"`
	TelcoName string `xml:"telcoName"`
}
