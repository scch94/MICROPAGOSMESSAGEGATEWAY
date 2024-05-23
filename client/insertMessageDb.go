package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

func CallToInsertMessageDB(validationStruct *helper.ToValidate, utfi string, ctx context.Context) helper.InsertResult {

	//traemos el contexto y le setiamos el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, "client")

	//creamos el struct para controlar la respuesta del smsgateway
	insertResult := helper.InsertResult{}

	//llenamos la estrucutra que se convertira en el json para enviar a la abse de datos
	ins_log.Infof(ctx, "PETITION[%v], startin to prepare the call to databasegateway to INSERT message", utfi)
	insertRequest := request.NewInsertMessageRequest(validationStruct, utfi)

	//preparamos el request pasando el struct que tiene el json a enviar la info
	req, err := prepareInsertMessageRequest(*insertRequest, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to prepareInsertMessageRequest()", utfi)
		insertResult.Id = ""
		insertResult.Message = "error when we try to prepareInsertMessageRequest()"
		insertResult.Result = err.Error()
		return insertResult
	}

	//hacemos el llamado y obtenemos el resultado
	insertMessageResponse, err := callToMicropagosInsertMessageDatabase(req, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to callToMicropagosDatabase", utfi)
		insertResult.Id = ""
		insertResult.Message = "error when we try to callToMicropagosDatabase()"
		insertResult.Result = err.Error()
		return insertResult
	} else {
		insertResult.Id = strconv.Itoa(insertMessageResponse.Id)
		insertResult.Message = insertMessageResponse.Message
		insertResult.Result = strconv.Itoa(insertMessageResponse.Result)
	}

	return insertResult
}
func prepareInsertMessageRequest(insertMessageRequest request.InsertMessageRequest, utfi string, ctx context.Context) (*http.Request, error) {

	//generamos el cuerpo de la solicitud para insertar el mensaje
	ins_log.Tracef(ctx, "PETITION[%v], starting to prepare the URL and the body of the petition", utfi)

	//agregamos lo params de la solicitud
	finalURL := fmt.Sprintf("%s/%s", config.Config.InsertMessage.URL, utfi)

	body, err := json.Marshal(insertMessageRequest)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to generate the body json.marshal()", utfi)
		return nil, err
	}

	//creamos la solicitud para ir a insetart el mensaje
	req, err := http.NewRequest(config.Config.InsertMessage.Method, finalURL, bytes.NewReader(body))
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error creating request to telcogateway: %v", err.Error(), utfi)
		return nil, err
	}
	ins_log.Tracef(ctx, "PETITION[%v], petition http created", utfi)

	//llenamos la cabecera
	req.Header.Set("Content-Type", "application/json")

	//logueamos la url y el body final de la peticion
	ins_log.Infof(ctx, "PETITION[%v], Final url : %s", utfi, finalURL)
	ins_log.Infof(ctx, "PETITION[%v], Final BODY: %s", utfi, body)
	return req, nil
}
func callToMicropagosInsertMessageDatabase(req *http.Request, utfi string, ctx context.Context) (response.InsertMessageResponse, error) {

	//creamos la variable que obtendra la respuesta de portabilidad
	var insertMessageResponse response.InsertMessageResponse

	//creamos el client, generamos el cronometro y realizamos peticion
	client := &http.Client{
		Timeout: time.Duration(config.Config.InsertMessage.Timeout) * time.Millisecond,
	}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error when we do the petition to micropagos databse: %s", utfi, err)
		return insertMessageResponse, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	ins_log.Infof(ctx, "PETITION[%v], Request to database tooks %v", utfi, duration)

	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		ins_log.Errorf(ctx, "PETITION[%v], error due to non-200 status code: %v", utfi, err)
		return insertMessageResponse, err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error reading response body: %s", utfi, err)
		return insertMessageResponse, err
	}

	// Imprimir la respuesta recibida
	statusCode := resp.StatusCode
	ins_log.Infof(ctx, "PETITION[%v], HTTP Status Response: %d", utfi, statusCode)
	ins_log.Infof(ctx, "PETITION[%v], RESPONSE BODY: %s", utfi, string(responseBody))
	//parceamos el resultado con lo que esperamos recibir

	//parceamos el resultado con lo que esperamos recibir
	err = json.Unmarshal(responseBody, &insertMessageResponse)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error decoding the response: %s", utfi, err)
		return insertMessageResponse, err
	}

	// Imprimir la respuesta
	ins_log.Infof(ctx, "PETITION[%v], this is the id of the message in the database: %d", utfi, insertMessageResponse.Id)
	return insertMessageResponse, nil
}
