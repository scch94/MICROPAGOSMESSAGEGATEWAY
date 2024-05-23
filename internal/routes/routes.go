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

func SetupRouter(ctx context.Context) *gin.Engine {

	// Agregar el valor "packageName" al contexto
	ctx = ins_log.SetPackageNameInContext(ctx, "routes")

	// create a new gin router and register the handlers
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Agregar middlewares global
	router.Use(middleware.GlobalMiddleware(ctx))
	router.Use(middleware.Authenticathormidldleware(ctx))
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

	//traemos el contexto y le setiamos el contexto actual
	ctx := c.Request.Context()
	ctx = ins_log.SetPackageNameInContext(ctx, "handler")

	ins_log.Errorf(ctx, "Route  not found: url: %v, method: %v", c.Request.RequestURI, c.Request.Method)
	c.JSON(http.StatusNotFound, nil)
}
