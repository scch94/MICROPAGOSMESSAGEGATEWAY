package main

import (
	"context"
	"os"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/server"
	"github.com/scch94/ins_log"
)

//lint:ignore SA1029 "Using built-in type string as key for context value intentionally"
var ctx = context.WithValue(context.Background(), "packageName", "main")

func main() {
	//levantamos la config
	errConfig := config.Upconfig()
	if errConfig != nil {
		ins_log.Errorf(ctx, "error when we try to get the configuration err: %v", errConfig)
		return
	}
	//inicialisamos el logger
	ins_log.StartLogger()
	ins_log.SetService("micropagosmessagegateway")
	ins_log.Infof(ctx, "startin micropagos message gateway version: %+v", version())
	//abrimos el archivo de logeo acutlamente solo logeara en consola
	file, err := os.OpenFile("logfile.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to get the log file : %s", err.Error())
		return
	}
	defer file.Close()
	//inicamos el servidor
	err = server.StartServer()
	if err != nil {
		ins_log.Errorf(ctx, "error al tratarde iniciar el servidor : %s", err.Error())
		return
	}

}

func version() string {
	return "1.0.0"
}
