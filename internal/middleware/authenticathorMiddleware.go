package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

// Authenticathormidldleware extrae y valida la autenticación básica del encabezado de la petición.
func Authenticathormidldleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ins_log.Infof(ctx, "starting to validate the ahutentication")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ins_log.Errorf(ctx, "missing authorization header")
			result := response.Result{
				DatabaseID: "",
				Message:    "Unauthorized",
				Code:       "12",
			}
			results := []response.Result{result}
			xmlresponse := response.GenerateXML(results)
			c.Data(http.StatusUnauthorized, "application/xml", []byte(xmlresponse))
			return
		}

		// Verifica que el esquema sea "Basic"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Basic {base64}"})
			result := response.Result{
				DatabaseID: "",
				Message:    "Unauthorized",
				Code:       "12",
			}
			results := []response.Result{result}
			xmlresponse := response.GenerateXML(results)
			c.Data(http.StatusUnauthorized, "application/xml", []byte(xmlresponse))
			return
		}

		// Decodifica el valor base64
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			ins_log.Errorf(ctx, "error decoding de value")
			result := response.Result{
				DatabaseID: "",
				Message:    "Unauthorized",
				Code:       "12",
			}
			results := []response.Result{result}
			xmlresponse := response.GenerateXML(results)
			c.Data(http.StatusUnauthorized, "application/xml", []byte(xmlresponse))
			return
		}

		// Separa el usuario y la contraseña
		userPass := strings.SplitN(string(decoded), ":", 2)
		if len(userPass) != 2 {
			ins_log.Errorf(ctx, "Invalid authorization format")
			result := response.Result{
				DatabaseID: "",
				Message:    "Unauthorized",
				Code:       "12",
			}
			results := []response.Result{result}
			xmlresponse := response.GenerateXML(results)
			c.Data(http.StatusUnauthorized, "application/xml", []byte(xmlresponse))
			return
		}

		username := userPass[0]
		password := userPass[1]
		ins_log.Infof(ctx, "starting Authentication proccess for %s", username)
		// Aquí validamos que el usuario y la contraseña sea la correcta por ahora quemamos resultado
		if username != "eaguerre" || password != "Mp3303" {
			ins_log.Errorf(ctx, "Invalid username or password")
			result := response.Result{
				DatabaseID: "",
				Message:    "Invalid username or password",
				Code:       "9",
			}
			results := []response.Result{result}
			xmlresponse := response.GenerateXML(results)
			c.Data(http.StatusUnauthorized, "application/xml", []byte(xmlresponse))
			return
		}
		ins_log.Infof(ctx, "the user:%s is already authenticated", username)

		//guardamos el nombre de usario en l contexto
		c.Set("username", username)

		// Si la validación es exitosa, continúa con el próximo handler
		c.Next()
	}
}
