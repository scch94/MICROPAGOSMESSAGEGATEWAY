package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/client"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/server"
	"github.com/scch94/ins_log"
)

func main() {

	// Creamos el contexto para esta ejecución
	ctx := context.Background()

	// Load configuration
	if err := config.Upconfig(ctx); err != nil {
		ins_log.Errorf(ctx, "error loading configuration: %v", err)
		return
	}

	go initializeAndWatchLogger(ctx)

	// Set logger configuration
	ins_log.SetService(serviceName)
	ins_log.SetLevel(config.Config.LogLevel)
	ctx = ins_log.SetPackageNameInContext(ctx, moduleName)
	ins_log.Infof(ctx, "starting micropagos message gateway version: %+v", getVersion())

	//inicamos el client
	client.InitHttpClient()

	// Start scheduled tasks
	startScheduler(ctx)

	// Start server
	go startServer(ctx)

	// Keep the program running
	select {}
}

func startServer(ctx context.Context) {
	if err := server.StartServer(ctx); err != nil {
		ins_log.Errorf(ctx, "error starting server: %s", err.Error())
	}
}

// esta funcion se encarga de traer el valor del mask cada que lo necesite: si el mask esta vacio y tiene un error al tratar de llamarlo vuelve a intentarlo cada minuto !
// esto garantiza que no importa el orden en el que levantes el modulo de la base y el modulo message gateway
func updateMask(ctx context.Context) {
	for {
		maskResponse, err := client.CallToGetMask(ctx)
		if err != nil {
			ins_log.Errorf(ctx, "error getting mask: %v", err)
			if len(helper.Mask) == 0 {
				ins_log.Warn(ctx, "mask is empty after error, retrying...")
				time.Sleep(1 * time.Minute)
				continue
			}
			break
		}
		helper.Mask = maskResponse.Masks
		ins_log.Trace(ctx, "mask was updated")
		break
	}
}

// proceso automatico que ira actualizando cada n tiempo el mask(el valor del timepo esta en la config)
func startScheduler(ctx context.Context) {
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.Every(config.Config.UpdateMaskTimeInMinutes).Minutes().Do(func() {
		updateMask(ctx)
	})
	go scheduler.StartAsync()
}

// funcion que ira cambiando de log cada hora
func initializeAndWatchLogger(ctx context.Context) {
	var file *os.File
	var logFileName string
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logDir := "../log"

			// Create the log directory if it doesn't exist
			if err = os.MkdirAll(logDir, 0755); err != nil {
				ins_log.Errorf(ctx, "error creating log directory: %v", err)
				return
			}

			// Define the log file name
			today := time.Now().Format("2006-01-02 15")
			replacer := strings.NewReplacer(" ", "_")
			today = replacer.Replace(today)
			logFileName = filepath.Join(logDir, config.Config.Log_name+today+".log")

			// Open the log file
			file, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				ins_log.Errorf(ctx, "error opening log file: %v", err)
				return
			}

			// Create a writer that writes to both file and console
			multiWriter := io.MultiWriter(os.Stdout, file)
			ins_log.StartLoggerWithWriter(multiWriter)

			// Esperar hasta el inicio de la próxima hora
			nextHour := time.Now().Truncate(time.Hour).Add(time.Hour)
			time.Sleep(time.Until(nextHour))

			// Close the previous log file
			file.Close()
		}
	}
}

func getVersion() string {
	return version
}
