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

// esta funcion al ser secundaria solo logueara errores !
func CallToGetMask(ctx context.Context) (response.MaskResponse, error) {

	//traemos el contexto y le setiamos el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, moduleName)

	//creamos el objeto que guardara la respuesta de la peticion
	masks := response.MaskResponse{}
	//ins_log.Tracef(ctx, "starting to prepare the call to databasegateway getmask")

	//preparamos el request
	req, err := prepareGetMaskReq(ctx)
	if err != nil {
		//ins_log.Errorf(ctx, "error when we try to prepareGetUserReq()")
		return masks, err
	}

	//hacemos el llamado y obtenemos el resultado
	masks, err = callToMicropagosGetMask(req, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to  callToMicropagosGetMask()")
		return masks, err
	}
	return masks, nil

}
func prepareGetMaskReq(ctx context.Context) (*http.Request, error) {

	//armaremos la url

	utfi := ins_log.GenerateUTFI()

	//agregamos lo params de la solicitud
	finalURL := fmt.Sprintf("%s/%s", config.Config.GetMask.URL, utfi)

	//creamos la solicitud para hacer el get filtered
	req, err := http.NewRequest(config.Config.GetMask.Method, finalURL, nil)
	if err != nil {
		ins_log.Errorf(ctx, "Error creating request to databasegateway: %v", err.Error())
		return nil, err
	}

	return req, nil

}

func callToMicropagosGetMask(req *http.Request, ctx context.Context) (response.MaskResponse, error) {

	var maskResponse response.MaskResponse

	resp, err := Client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "Error when we do the petition to micropagos databse: %s", err)
		return maskResponse, err
	}
	defer resp.Body.Close()

	//confirmamos que la respueta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		return maskResponse, err
	}

	//logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "Error reading response body: %s", err)
		return maskResponse, err
	}

	//parceamos el resultado con lo que esperamos recibir
	err = json.Unmarshal(responseBody, &maskResponse)
	if err != nil {
		return maskResponse, err
	}

	return maskResponse, nil
}
