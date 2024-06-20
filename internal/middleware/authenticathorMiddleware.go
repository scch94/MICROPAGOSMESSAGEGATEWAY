package middleware

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/client"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

// Authenticathormidldleware extrae y valida la autenticación básica del encabezado de la petición.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//traemos el contexto y le setiamos el contexto actual
		ctx := c.Request.Context()
		ctx = ins_log.SetPackageNameInContext(ctx, "middleware")
		ins_log.Tracef(ctx, "starting to validate the ahutentication")

		username, password, ok := c.Request.BasicAuth()

		if !ok {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			sendErrorResponse(ctx, c, "missing authorization header", constants.ERROR_UNAUTHORIZED)
			return
		}
		// if authHeader == "" {
		// 	ins_log.Errorf(ctx, "missing authorization header")
		// 	sendErrorResponse(ctx, c, "missing authorization header", constants.ERROR_UNAUTHORIZED)
		// 	return
		// }

		// Aquí validamos que el usuario y la contraseña sea la correcta por ahora quemamos resultado
		userData, err := helper.GetUserdata(ctx, username)
		if err != nil {
			ins_log.Errorf(ctx, "error getting user data to the user %s,User does not exist or does not have permission to use the service.", username)
			sendErrorResponse(ctx, c, "User does not exist or does not have permission to use the service.", constants.ERROR_INVALID_USERNAME_OR_PASSWORD)
			return
		}

		//formateamos la password
		formatedPassword := formatPassword(password)
		ins_log.Tracef(ctx, "formatted password: %s", formatedPassword)

		if username != userData.Username || formatedPassword != userData.Password {
			sendErrorResponse(ctx, c, "Invalid password", constants.ERROR_INVALID_USERNAME_OR_PASSWORD)
			return
		}

		ins_log.Infof(ctx, "the user:%s is already authenticated", username)

		//guardamos el nombre de usarrio en el contexto y tambien el dominio que luego utilizaremos
		c.Set("username", userData.Username)
		c.Set("dominio", userData.UserDomain)

		//proceso que insertara el valor del last login del usuario
		go startUpdatelastLogin(ctx, userData.Username)

		// Si la validación es exitosa, continúa con el próximo handler
		c.Next()
	}
}

// formateador de password
func formatPassword(password string) string {
	// Concatenar la salt y la contraseña en un slice de bytes.
	saltedPassword := append(constants.SALT, []byte(password)...)

	// Calcular el hash SHA-1 de la combinación de sal y contraseña.
	hash := sha1.New()
	hash.Write(saltedPassword)
	hashedBytes := hash.Sum(nil)

	// Codificar el hash resultante en Base64.
	hashedPassword := base64.StdEncoding.EncodeToString(hashedBytes)

	return hashedPassword
}

// GENERA RESPUESTA
func sendErrorResponse(ctx context.Context, c *gin.Context, message, code string) {
	ins_log.Errorf(ctx, message)
	result := response.Result{
		DatabaseID: "",
		Message:    message,
		Code:       code,
	}
	results := []response.Result{result}
	xmlResponse := response.GenerateXML(results)
	c.Data(http.StatusUnauthorized, "application/xml", []byte(xmlResponse))
	c.Abort()
}

func startUpdatelastLogin(ctx context.Context, username string) {

	// Get the current time
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")

	UserLastLogin := request.UserData{
		UserName:  username,
		LoginTime: formattedTime,
	}

	userToUpdate := request.UsersToUpdate{
		Users: []request.UserData{
			UserLastLogin,
		},
	}

	_, err := client.CallToUpdateLastLogin(ctx, userToUpdate)
	if err != nil {
		ins_log.Errorf(ctx, "error updating last login: %v", err)
		return
	}
	ins_log.Infof(ctx, "the last login of the user: %s was updated the new value is: %s", username, formattedTime)
}
