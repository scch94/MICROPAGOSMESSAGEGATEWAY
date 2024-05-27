package main

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/server"
	"github.com/scch94/ins_log"
)

func main() {

	// Creamos el contexto para esta ejecución
	ctx := context.Background()

	today := time.Now().Format("2006-01-02 15")
	// Reemplazar los caracteres no permitidos en el nombre del archivo
	replacer := strings.NewReplacer(" ", "_")
	today = replacer.Replace(today)

	// Construir el nombre del archivo de log
	logFileName := "micropagosmessagegateway_" + today + ".log"
	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Creamos un escritor que escriba tanto en el archivo como en la consola
	multiWriter := io.MultiWriter(os.Stdout, file)
	ins_log.StartLoggerWithWriter(multiWriter)

	// Levantamos la configuración
	errConfig := config.Upconfig(ctx)
	if errConfig != nil {
		ins_log.Errorf(ctx, "error when we try to get the configuration err: %v", errConfig)
		return
	}

	// Inicializamos el logger
	ins_log.SetService("micropagosmessagegateway")
	ins_log.SetLevel(config.Config.LogLevel)

	// Agregamos el valor "packageName" al contexto
	ctx = ins_log.SetPackageNameInContext(ctx, "main")

	ins_log.Infof(ctx, "starting micropagos message gateway version: %+v", version())

	// Iniciamos el servidor
	err = server.StartServer(ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error al tratar de iniciar el servidor: %s", err.Error())
		return
	}
}

func version() string {
	return "1.0.0"
}
