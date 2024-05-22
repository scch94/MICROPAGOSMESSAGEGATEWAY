package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

func CallToGetUserDomain(validationStruct *helper.ToValidate, utfi string, ctx context.Context) helper.UserDomainResult {

	//traemos el contexto y le setiamos el contexto actual
	ctx = context.WithValue(ctx, constants.PACKAGE_NAME_KEY, "client")

	//cramos el struct que nos ayudara a contraolar la respuesta del GETUSERDOMAIN en la base de datos
	getUserDomainResult := helper.UserDomainResult{}

	//creamos la estrucutra que tendra los params
	ins_log.Infof(ctx, "PETITION[%v], startin to prepare the call to databasegateway GetUserDomain", utfi)
	params := request.NewGetUserDomainRequest(validationStruct.Username)

	//preparamos el request pasando los params de la peticion para insertarlos en la
	req, err := prepareGetUserReq(*params, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to prepareGetUserReq()", utfi)
		getUserDomainResult.UserDomainError = err
		getUserDomainResult.UserDomainMessage = "error when we try to prepareGetUserReq"
		getUserDomainResult.UserDomainResult = err.Error()
		return getUserDomainResult
	}

	//hacemos el llamado y obtenemos el resultado
	getUserDomainresponse, err := callToMicropagosGetUserDomainDatabase(req, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to  callToMicropagosGetUserDomainDatabase()", utfi)
		getUserDomainResult.UserDomainError = err
		getUserDomainResult.UserDomainMessage = "error when we try to prepareGetUserReq"
		getUserDomainResult.UserDomainResult = err.Error()
		return getUserDomainResult
	} else {
		getUserDomainResult.UserDomainResult = getUserDomainresponse.UserDomain
		getUserDomainResult.UserDomainMessage = getUserDomainresponse.Message
		getUserDomainResult.UserDomainError = nil
		return getUserDomainResult
	}

}

func prepareGetUserReq(params request.GetUserDomainRequest, utfi string, ctx context.Context) (*http.Request, error) {

	//armaremos la url para agregar los params a la url
	ins_log.Tracef(ctx, "PETITION[%v], starting to prepare the URL to the petition", utfi)

	//agregamos lo params de la solicitud
	finalURL := fmt.Sprintf("%s/%s/%s", config.Config.GetUserDomain.URL, params.UserName, utfi)

	//creamos la solicitud para hacer el get filtered
	req, err := http.NewRequest(config.Config.GetUserDomain.Method, finalURL, nil)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error creating request to databasegateway: %v", err.Error(), utfi)
		return nil, err
	}
	ins_log.Tracef(ctx, "PETITION[%v], petition http created", utfi)

	//logeamos la url
	ins_log.Infof(ctx, "PETITION[%v], Final url : %s", utfi, finalURL)

	return req, nil
}

func callToMicropagosGetUserDomainDatabase(req *http.Request, utfi string, ctx context.Context) (response.UserDomainResponse, error) {

	var UserDomainResponse response.UserDomainResponse

	//creamos el client, generamos el cronometro y realizamos la peticion
	client := &http.Client{
		Timeout: time.Duration(config.Config.GetUserDomain.Timeout) * time.Millisecond,
	}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error when we do the petition to micropagos databse: %s", utfi, err)
		return UserDomainResponse, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	ins_log.Infof(ctx, "PETITION[%v], Request to Micropagos_database tooks %v", utfi, duration)

	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		ins_log.Errorf(ctx, "PETITION[%v], error due to non-200 status code: %v", utfi, err)
		return UserDomainResponse, err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error reading response body: %s", utfi, err)
		return UserDomainResponse, err
	}

	statusCode := resp.StatusCode
	ins_log.Infof(ctx, "PETITION[%v], HTTP Status Response: %d", utfi, statusCode)
	ins_log.Infof(ctx, "PETITION[%v], RESPONSE BODY: %s", utfi, string(responseBody))

	//parceamos el resultado con lo que esperamos recibir
	err = json.Unmarshal(responseBody, &UserDomainResponse)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error decoding the response: %s", utfi, err)
		return UserDomainResponse, err
	}

	ins_log.Debugf(ctx, "the user domain is: %v", UserDomainResponse.UserDomain)

	return UserDomainResponse, nil

}
