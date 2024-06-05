package main

import (
	"context"
	"io"
	"os"
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

	logFileName := initializeLogger()
	defer logFileName.Close()

	// Load configuration
	if err := config.Upconfig(ctx); err != nil {
		ins_log.Errorf(ctx, "error loading configuration: %v", err)
		return
	}

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
func initializeLogger() *os.File {
	today := time.Now().Format("2006-01-02_15")
	logFileName := logFileName + today + ".log"
	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	ins_log.StartLoggerWithWriter(multiWriter)
	return file
}

func getVersion() string {
	return version
}

func startServer(ctx context.Context) {
	if err := server.StartServer(ctx); err != nil {
		ins_log.Errorf(ctx, "error starting server: %s", err.Error())
	}
}
func updateMask(ctx context.Context) {
	maskResponse, err := client.CallToGetMask(ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error gettin mask")
		return
	}
	helper.Mask = maskResponse.Masks
	ins_log.Trace(ctx, "mask was updated")
}
func startScheduler(ctx context.Context) {
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.Every(config.Config.UpdateMaskTimeInMinutes).Minutes().Do(func() {
		updateMask(ctx)
	})
	go scheduler.StartAsync()
}
