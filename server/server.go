package server

import (
	"context"
	"net/http"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/routes"
	"github.com/scch94/ins_log"
)

func StartServer(ctx context.Context) error {
	// actualizamos contexto y logueamos el puerto
	ctx = context.WithValue(ctx, constants.PACKAGE_NAME_KEY, "server")
	ins_log.Infof(ctx, "Starting server on address: %s", config.Config.ServPort)
	//usamos las rutas
	router := routes.SetupRouter(ctx)
	serverConfig := &http.Server{
		Addr:    config.Config.ServPort,
		Handler: router,
	}
	err := serverConfig.ListenAndServe()
	if err != nil {
		ins_log.Errorf(ctx, "cant connect to the server: %+v", err)
		return err
	}
	return nil
}
