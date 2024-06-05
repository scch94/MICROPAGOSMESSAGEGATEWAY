package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

func CallToGetUserData(request request.GetUserRequest, ctx context.Context) (response.UserResponse, error) {

	//traemos el contexto y le setiamos el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, moduleName)

	//creamos el objeto que guardara la respuesta de la peticion
	user := response.UserResponse{}
	ins_log.Infof(ctx, "startin to prepare the call to databasegateway GetUserDomain")

	//preparamos el request pasando los params de la peticion para insertarlos en la url
	req, err := prepareGetUserReq(request, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to prepareGetUserReq()")
		return user, err
	}

	//hacemos el llamado y obtenemos el resultado
	user, err = callToMicropagosGetUserDatabase(req, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to  callToMicropagosGetUserDomainDatabase()")
		return user, err
	}
	return user, nil

}

func prepareGetUserReq(params request.GetUserRequest, ctx context.Context) (*http.Request, error) {

	//armaremos la url para agregar los params a la url
	ins_log.Tracef(ctx, "starting to prepare the URL to the petition")

	utfi := ins_log.GetUTFIFromContext(ctx)

	//agregamos lo params de la solicitud
	finalURL := fmt.Sprintf("%s/%s/%s", config.Config.GetUserDomain.URL, params.UserName, utfi)

	//creamos la solicitud para hacer el get filtered
	req, err := http.NewRequest(config.Config.GetUserDomain.Method, finalURL, nil)
	if err != nil {
		ins_log.Errorf(ctx, "Error creating request to databasegateway: %v", err.Error())
		return nil, err
	}
	ins_log.Tracef(ctx, "petition http created")

	//logeamos la url
	ins_log.Infof(ctx, "Final url : %s", finalURL)

	return req, nil
}

func callToMicropagosGetUserDatabase(req *http.Request, ctx context.Context) (response.UserResponse, error) {

	var UserResponse response.UserResponse

	start := time.Now()

	resp, err := Client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "Error when we do the petition to micropagos databse: %s", err)
		return UserResponse, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	ins_log.Infof(ctx, "Request to Micropagos_database tooks %v", duration)

	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		ins_log.Errorf(ctx, "error due to non-200 status code: %v", err)
		return UserResponse, err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "Error reading response body: %s", err)
		return UserResponse, err
	}

	statusCode := resp.StatusCode
	ins_log.Infof(ctx, "HTTP Status Response: %d", statusCode)
	ins_log.Infof(ctx, "RESPONSE BODY: %s", string(responseBody))

	//parceamos el resultado con lo que esperamos recibir
	err = json.Unmarshal(responseBody, &UserResponse)
	if err != nil {
		ins_log.Errorf(ctx, "Error decoding the response: %s", err)
		return UserResponse, err
	}

	ins_log.Tracef(ctx, "the getUser response is: %v", UserResponse)

	return UserResponse, nil

}
