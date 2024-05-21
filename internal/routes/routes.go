package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/handler"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/middleware"
	"github.com/scch94/ins_log"
)

//lint:ignore SA1029 "Using built-in type string as key for context value intentionally"
var ctx = context.WithValue(context.Background(), "packageName", "routes")

func SetupRouter() *gin.Engine {
	// create a new gin router and register the handlers
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Agregar middlewares global
	router.Use(middleware.GlobalMiddleware())
	router.Use(middleware.Authenticathormidldleware())
	router.Use(gin.Recovery())

	h := handler.NewHandler()

	//metodos
	router.GET("/", h.Welcome)
	router.POST(constants.PATH, h.SendMessageService)
	router.NoRoute(notFoundHandler)
	return router
}

// Controlador para manejar rutas no encontradas
func notFoundHandler(c *gin.Context) {

	ins_log.Errorf(ctx, "Route  not found: url: %v, method: %v", c.Request.RequestURI, c.Request.Method)
	c.JSON(http.StatusNotFound, nil)
}
