package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

func CallToGetUsers(ctx context.Context) (response.GetUsersInfoResponse, error) {

	//traemos el contexto y le setiamos el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, moduleName)

	//creamos el objeto que guardara la respuesta de la peticion
	Users := response.GetUsersInfoResponse{}

	//preparamos el request
	req, err := prepareGetUsersReq(ctx)
	if err != nil {
		//ins_log.Errorf(ctx, "error when we try to prepareGetUserReq()")
		return Users, err
	}

	//hacemos el llamado y obtenemos el resultado
	Users, err = callToMicropagosGetUsersInfo(req, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to  callToMicropagosGetUsersInfo()")
		return Users, err
	}
	return Users, nil
}

func prepareGetUsersReq(ctx context.Context) (*http.Request, error) {

	//armamos la url
	utfi := ins_log.GenerateUTFI()

	//agregamos lo params de la solicitud
	finalURL := fmt.Sprintf("%s/%s", config.Config.GetUsersInfo.URL, utfi)

	//creamos la solicitud para hacer el get filtered
	req, err := http.NewRequest(config.Config.GetUsersInfo.Method, finalURL, nil)
	if err != nil {
		ins_log.Errorf(ctx, "Error creating request to databasegateway: %v", err.Error())
		return nil, err
	}

	return req, nil

}

func callToMicropagosGetUsersInfo(req *http.Request, ctx context.Context) (response.GetUsersInfoResponse, error) {

	var usersInfoResponse response.GetUsersInfoResponse

	resp, err := Client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "Error when we do the petition to micropagos database: %v", err.Error())
		return usersInfoResponse, err
	}
	defer resp.Body.Close()
	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		return usersInfoResponse, err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "Error reading response body: %s", err)
		return usersInfoResponse, err
	}

	//parceamos el resultado con lo que esperamos recibir
	err = json.Unmarshal(responseBody, &usersInfoResponse)
	if err != nil {
		return usersInfoResponse, err
	}
	return usersInfoResponse, nil
}
