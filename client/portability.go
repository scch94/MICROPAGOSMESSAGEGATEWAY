package client

//"url_smsgateway": "http://192.168.27.41:8090/sendsms",
import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

func CallPortabilidad(validationStruct *helper.ToValidate, utfi string, ctx context.Context) helper.PortabilidadResult {

	//traemos el contexto y le setiamos el contexto actual
	ctx = context.WithValue(ctx, constants.PACKAGE_NAME_KEY, "client")

	//creamos el struct para controlar la respuesta de portabilidad
	portabilidadResult := helper.PortabilidadResult{}

	//vamos a prepareRequest para llenar el request antes de enviarlo
	ins_log.Infof(ctx, "PETITION[%v], startin to prepare the call to portability whit number %s", utfi, validationStruct.Mobile)
	req, err := prepareRequest(validationStruct.Mobile, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to prepareRequest()", utfi)
		portabilidadResult.PassedPortabilidad = false
		portabilidadResult.PortabilidadMessage = err.Error()
		return portabilidadResult
	}

	//hacemos el llamado y ontenemos el nombre de la telco, lo devolvemos directo por que la funcion ya devuelve la palabra y un error
	validationStruct.Telco, err = callToPortabilidad(req, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to call to portabiliad()", utfi)
		portabilidadResult.PassedPortabilidad = false
		portabilidadResult.PortabilidadMessage = err.Error()
		return portabilidadResult
	}
	portabilidadResult.PassedPortabilidad = true
	portabilidadResult.PortabilidadMessage = "OK"
	return portabilidadResult

}

func prepareRequest(msisdn string, utfi string, ctx context.Context) (*http.Request, error) {

	//generamos el cuerpo de la solicitud de la portabilidad
	bodyToPortabildiad, err := request.CreateBodyToPortabilidad(msisdn)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], problem when we try to create body to portability error: %v", utfi, err.Error())
		return nil, err
	}

	//creamos la solicitud para ir a portabilidad
	req, err := http.NewRequest(config.Config.Portabilidad.Method, config.Config.Portabilidad.URL, strings.NewReader(bodyToPortabildiad))
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error creating request to portability: %v", utfi, err.Error())
		return nil, err
	}
	ins_log.Tracef(ctx, "PETITION[%v], petition http created", utfi)

	//llenamos cabecera y basic authenticator
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", `"urn:getTelco"`)
	req.SetBasicAuth(config.Config.Portabilidad.Username, config.Config.Portabilidad.Password)

	//logueamos la url y el body final de la peticion
	ins_log.Infof(ctx, "PETITION[%v], Final url : %s", utfi, config.Config.Portabilidad.URL)
	ins_log.Infof(ctx, "PETITION[%v], Final BODY: %s", utfi, bodyToPortabildiad)
	return req, nil
}

func callToPortabilidad(req *http.Request, utfi string, ctx context.Context) (string, error) {

	//creamos la variable que obtendra la respuesta de portabilidad
	var portabiliadResponse response.PortabilidadResponse

	//creamos el client generamos cronometro y realizamos la peticion
	client := &http.Client{
		Timeout: time.Duration(config.Config.Portabilidad.Timeout) * time.Millisecond,
	}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error when we do the petition to portabilidad: %s", utfi, err)
		return "", err
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	ins_log.Infof(ctx, "PETITION[%v], Request to PORTABILIDAD tooks %v", utfi, duration)

	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		ins_log.Errorf(ctx, "PETITION[%v], error due to non-200 status code: %v", utfi, err)
		return "", err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error reading response body: %s", utfi, err)
		return "", err
	}

	// Imprimir la respuesta recibida
	statusCode := resp.StatusCode
	ins_log.Infof(ctx, "PETITION[%v], HTTP Status Response: %d", utfi, statusCode)
	ins_log.Infof(ctx, "PETITION[%v], RESPONSE BODY: %s", utfi, string(responseBody))

	//parceamos el resultado con lo que esperamos recibir
	err = xml.Unmarshal(responseBody, &portabiliadResponse)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error decoding the response: %s", utfi, err)
		return "", err
	}

	// check if the responses is unkown
	if portabiliadResponse.Body.GetTelcoResponse.Return.TelcoName == "UNKNOWN" {
		ins_log.Errorf(ctx, "PETITION[%v], the number is not vinculated to any telco number", utfi)
		return "", errors.New("error trying to get telco")
	}
	return portabiliadResponse.Body.GetTelcoResponse.Return.TelcoName, err
}
