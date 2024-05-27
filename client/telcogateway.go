package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

func CallTelcoGateway(validationStruct *helper.ToValidate, utfi string, ctx context.Context) helper.SmsgatewayResult {

	//traemos el contexto y le setiamos el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, "client")

	//creamos el struct para controlar la respuesta del smsgateway
	telcoGatewayResullt := helper.SmsgatewayResult{}

	//creamos el cuerpo de la peticion
	ins_log.Infof(ctx, "PETITION[%v], startin to prepare the call to telcoGateway", utfi)
	smsGatewayRequest := request.NewSmsGatewayRequest(*validationStruct, utfi)
	smsGatewayRequest.TLVValue = validationStruct.Telco
	smsGatewayRequest.TLVLength = len(validationStruct.Telco)

	//preparamos el request pasando el struct que tiene el json a enviar la info
	req, err := prepareTelcoGatewayRequest(*smsGatewayRequest, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to prepareTelcoGatewayRequest()", utfi)
		telcoGatewayResullt.PassedSmsgateway = false
		telcoGatewayResullt.SmsgatewayDescription = "error when we try to prepareTelcoGatewayRequest()"
		telcoGatewayResullt.SmsgatewayResult = err.Error()
		return telcoGatewayResullt
	}

	//hacemos el llamado y obtenemos el resultado
	responseGateway, err := calltoTelcoGatewayRequest(req, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to call to calltoTelcoGatewayRequest()", utfi)
		telcoGatewayResullt.PassedSmsgateway = false
		telcoGatewayResullt.SmsgatewayDescription = "error when we try to calltoTelcoGatewayRequest()"
		telcoGatewayResullt.SmsgatewayResult = err.Error()
		return telcoGatewayResullt
	}
	telcoGatewayResullt.PassedSmsgateway = true
	telcoGatewayResullt.SmsgatewayDescription = responseGateway.Description
	telcoGatewayResullt.SmsgatewayResult = responseGateway.Status
	return telcoGatewayResullt
	// return calltoTelcoGatewayRequest(req)
}
func prepareTelcoGatewayRequest(smsGatewayRequest request.SmsGatewayRequest, utfi string, ctx context.Context) (*http.Request, error) {

	//generamos el cuerpo de la solicitud para el gateway que enviara el mensaje
	ins_log.Tracef(ctx, "PETITION[%v], starting to prepare the body of the petition", utfi)
	body, err := json.Marshal(smsGatewayRequest)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to generate the body json.marshal()", utfi)
		return nil, err
	}

	//creamos la solicitud para ir a insetart el mensaje
	req, err := http.NewRequest(config.Config.SMSGateway.Method, config.Config.SMSGateway.URL, bytes.NewReader(body))
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error creating request to telcogateway: %v", utfi, err.Error())
		return nil, err
	}
	ins_log.Tracef(ctx, "PETITION[%v], petition http created", utfi)

	//llenamos cabecera
	req.Header.Set("Content-Type", "application/json")

	//logueamos la url y el body final de la peticion
	ins_log.Infof(ctx, "PETITION[%v], Final url : %s", utfi, config.Config.SMSGateway.URL)
	ins_log.Infof(ctx, "PETITION[%v], Final BODY: %s", utfi, body)
	return req, nil
}

func calltoTelcoGatewayRequest(req *http.Request, utfi string, ctx context.Context) (response.SmsGatewayResponse, error) {

	//creamos la variable que obtendra la respuesta de portabilidad
	var smsGatewayResponse response.SmsGatewayResponse

	//traemos el client y le configuramos el timeout , generamos el cronometro y realizamos la peticion
	client.Timeout = time.Duration(config.Config.SMSGateway.Timeout) * time.Millisecond
	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error when we do the petition to smsgateway: %s", utfi, err)
		return smsGatewayResponse, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	ins_log.Infof(ctx, "PETITION[%v], Request to SMSGATEWAY tooks %v", utfi, duration)

	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		ins_log.Errorf(ctx, "PETITION[%v], error due to non-200 status code: %v", utfi, err)
		return smsGatewayResponse, err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error reading response body: %s", utfi, err)
		return smsGatewayResponse, err
	}

	// Imprimir la respuesta recibida
	statusCode := resp.StatusCode
	ins_log.Infof(ctx, "PETITION[%v], HTTP Status Response: %d", utfi, statusCode)
	ins_log.Infof(ctx, "PETITION[%v], RESPONSE BODY: %s", utfi, string(responseBody))
	//parceamos el resultado con lo que esperamos recibir

	//parceamos el resultado con lo que esperamos recibir
	err = json.Unmarshal(responseBody, &smsGatewayResponse)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error decoding the response: %s", utfi, err)
		return smsGatewayResponse, err
	}

	// Imprimir la respuesta
	ins_log.Infof(ctx, "PETITION[%v], this is the code of the response: %s", utfi, smsGatewayResponse.Status)
	return smsGatewayResponse, nil
}
