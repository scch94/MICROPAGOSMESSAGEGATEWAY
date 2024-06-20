package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

func CallToUpdateLastLogin(ctx context.Context, userToUpdate request.UsersToUpdate) (response.UpdateUsersLastLoginResponse, error) {
	//traemos el contexto y le setiamos el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, moduleName)

	//creamos el objeto que guardara la respuesta de la peticion
	updateUsersLastLoginResponse := response.UpdateUsersLastLoginResponse{}

	//preparamos el request
	req, err := prepareUpdateLastLoginReq(ctx, userToUpdate)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to prepareUpdateLastLoginReq()")
		return updateUsersLastLoginResponse, err
	}

	//hacemos el llamado y obtenemos el resultado
	updateUsersLastLoginResponse, err = callToMicroDbUpdateLastLoginReq(req, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to  callToMicropagosGetUsersInfo()")
		return updateUsersLastLoginResponse, err
	}
	return updateUsersLastLoginResponse, nil

}

func prepareUpdateLastLoginReq(ctx context.Context, userToUpdate request.UsersToUpdate) (*http.Request, error) {

	//generamos el cuerpo de la solicitud para insertar el mensaje
	ins_log.Tracef(ctx, "starting to prepare the URL and the body of the petition")

	//armamos la url
	utfi := ins_log.GetUTFIFromContext(ctx)

	//agregamos lo params de la solicitud
	finalURL := fmt.Sprintf("%s/%s", config.Config.UpdateUsersLastLogin.URL, utfi)
	//ajuntamos el body de la peticion
	body, err := json.Marshal(userToUpdate)
	if err != nil {
		ins_log.Errorf(ctx, "[%v], error when we try to generate the body json.marshal()", utfi)
		return nil, err
	}

	//creamos la solicitud para hacer el get filtered
	req, err := http.NewRequest(config.Config.UpdateUsersLastLogin.Method, finalURL, bytes.NewReader(body))
	if err != nil {
		ins_log.Errorf(ctx, "Error creating request to databasegateway: %v", err.Error())
		return nil, err
	}
	return req, nil

}

func callToMicroDbUpdateLastLoginReq(req *http.Request, ctx context.Context) (response.UpdateUsersLastLoginResponse, error) {

	var updateUsersLastLoginResponse response.UpdateUsersLastLoginResponse

	resp, err := Client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "Error when we do the petition to micropagos database: %v", err.Error())
		return updateUsersLastLoginResponse, err
	}

	defer resp.Body.Close()

	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		return updateUsersLastLoginResponse, err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "Error reading response body: %s", err)
		return updateUsersLastLoginResponse, err
	}

	//parceamos el resultado con lo que esperamos recibir
	err = json.Unmarshal(responseBody, &updateUsersLastLoginResponse)
	if err != nil {
		ins_log.Errorf(ctx, "Error decoding the response: %s", err)
		return updateUsersLastLoginResponse, err
	}
	return updateUsersLastLoginResponse, nil
}
