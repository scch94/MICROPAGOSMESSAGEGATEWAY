package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scch94/ins_log"
)

func GlobalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generar un UTFI y agregarlo al contexto
		ctx := c.Request.Context()
		ctx = ins_log.SetPackageNameInContext(ctx, "middleware")
		utfi := ins_log.GenerateUTFI()
		ctx = ins_log.SetUTFIInContext(ctx, utfi)

		// Copiar el contexto a la solicitud
		c.Request = c.Request.WithContext(ctx)

		ins_log.Info(ctx, "New petition received")
		ins_log.Tracef(ctx, "url: %v, method: %v", c.Request.RequestURI, c.Request.Method)
		startTime := time.Now() // Registro de inicio de tiempo

		// Pasar la solicitud al siguiente middleware o al controlador final
		c.Next()

		//logeamos el tiempo final de la peticion
		elapsedTime := time.Since(startTime)
		ins_log.Infof(ctx, "Request took %v", elapsedTime)
	}
}
