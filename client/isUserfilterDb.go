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

func CallToFiltersDB(validationStruct *helper.ToValidate, utfi string, ctx context.Context) helper.FilterResult {

	//traemos el contexto y le setiamos el contexto actual
	ctx = context.WithValue(ctx, constants.PACKAGE_NAME_KEY, "client")

	//cramos el struct que nos ayudara a contraolar la respuesta del filter en la base de datos
	filterResult := helper.FilterResult{}

	//creamos la estrucutra que tendra los params
	ins_log.Infof(ctx, "PETITION[%v], startin to prepare the call to databasegateway checking FILTER", utfi)
	params := request.NewFilterRequest(validationStruct.Mobile, validationStruct.ShortNumber)

	//preparamos el request pasando los params de la peticion para insertarlos en la
	req, err := prepareFilterReq(*params, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to prepareInsertMessageRequest()", utfi)
		filterResult.IsFilter = false
		filterResult.FilterMessage = "errro when we try to prepare filterDB"
		filterResult.Error = err
		return filterResult
	}

	//hacemos el llamado y obtenemos el resultado
	filterMessageResponse, err := callToMicropagosFilterDatabase(req, utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error when we try to callToMicropagosFilterDatabase()", utfi)
		filterResult.IsFilter = false
		filterResult.FilterMessage = "errro when we try to call filterDB"
		filterResult.Error = err
		return filterResult
	} else {
		filterResult.IsFilter = filterMessageResponse.IsFilter
		filterResult.FilterMessage = filterMessageResponse.Message
		filterResult.Error = nil
		return filterResult
	}

}

func prepareFilterReq(params request.FilterRequest, utfi string, ctx context.Context) (*http.Request, error) {

	//armaremos la url para agregar los params a la url
	ins_log.Tracef(ctx, "PETITION[%v], starting to prepare the URL to the petition", utfi)

	//agregamos lo params de la solicitud
	finalURL := fmt.Sprintf("%s/%s/%s/%s", config.Config.GetFilterDatabase.URL, params.MobileNumber[3:], params.ShortNUmber, utfi)

	//creamos la solicitud para hacer el get filtered
	req, err := http.NewRequest(config.Config.GetFilterDatabase.Method, finalURL, nil)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error creating request to databasegateway: %v", err.Error(), utfi)
		return nil, err
	}
	ins_log.Tracef(ctx, "PETITION[%v], petition http created", utfi)

	//logeamos la url
	ins_log.Infof(ctx, "PETITION[%v], Final url : %s", utfi, finalURL)

	return req, nil
}

func callToMicropagosFilterDatabase(req *http.Request, utfi string, ctx context.Context) (response.FilterResponse, error) {

	var filterMessageResponse response.FilterResponse

	//creamos el client, generamos el cronometro y realizamos la peticion
	client := &http.Client{
		Timeout: time.Duration(config.Config.GetFilterDatabase.Timeout) * time.Second,
	}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error when we do the petition to micropagos databse: %s", utfi, err)
		return filterMessageResponse, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	ins_log.Infof(ctx, "PETITION[%v], Request to database tooks %v", utfi, duration)

	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		ins_log.Errorf(ctx, "PETITION[%v], error due to non-200 status code: %v", utfi, err)
		return filterMessageResponse, err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error reading response body: %s", utfi, err)
		return filterMessageResponse, err
	}

	// Imprimir la respuesta recibida
	statusCode := resp.StatusCode
	ins_log.Infof(ctx, "PETITION[%v], HTTP Status Response: %d", utfi, statusCode)
	ins_log.Infof(ctx, "PETITION[%v], RESPONSE BODY: %s", utfi, string(responseBody))

	//parceamos el resultado con lo que esperamos recibir
	err = json.Unmarshal(responseBody, &filterMessageResponse)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], Error decoding the response: %s", utfi, err)
		return filterMessageResponse, err
	}

	ins_log.Debugf(ctx, "the filter message response was: shortnumber and destinity are filter? : %v", filterMessageResponse.IsFilter)

	return filterMessageResponse, nil

}
