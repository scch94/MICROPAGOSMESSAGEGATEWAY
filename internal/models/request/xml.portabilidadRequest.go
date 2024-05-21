package request

import (
	"encoding/xml"
)

type CallPortabilidad struct {
	XMLName xml.Name         `xml:"Envelope"`
	Soapenv string           `xml:"xmlns:soapenv,attr"`
	Web     string           `xml:"xmlns:web,attr"`
	Header  string           `xml:"soapenv:Header"`
	Body    BodyPortabilidad `xml:"soapenv:Body"`
}

type BodyPortabilidad struct {
	GetTelco GetTelco `xml:"web:getTelco"`
}

type GetTelco struct {
	Msisdn string `xml:"web:msisdn"`
}

func NewEnvelopeFromXML(msisdn string) (*CallPortabilidad, error) {
	var envelope CallPortabilidad
	err := xml.Unmarshal([]byte(msisdn), &envelope)
	if err != nil {
		return nil, err
	}
	return &envelope, nil
}

// CreateSoapEnvelope crea un XML de solicitud SOAP con el número de teléfono proporcionado.
func CreateBodyToPortabilidad(msisdn string) (string, error) {

	envelope := CallPortabilidad{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Web:     "http://webservices.push.inswitch.us",
		Body: BodyPortabilidad{
			GetTelco: GetTelco{
				Msisdn: msisdn,
			},
		},
	}

	xmlData, err := xml.MarshalIndent(envelope, "", "\t")
	if err != nil {
		return "", err
	}

	return string(xmlData), nil
}
