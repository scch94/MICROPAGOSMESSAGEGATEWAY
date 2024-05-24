package middleware

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/client"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
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
		ins_log.Infof(ctx, "starting to validate the ahutentication")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ins_log.Errorf(ctx, "missing authorization header")
			sendErrorResponse(ctx, c, "missing authorization header", constants.ERROR_UNAUTHORIZED)
			return
		}

		// Verifica que el esquema sea "Basic"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			sendErrorResponse(ctx, c, "missing authorization header", constants.ERROR_UNAUTHORIZED)
			return
		}

		// Decodifica el valor base64
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			ins_log.Errorf(ctx, "error decoding de value")
			sendErrorResponse(ctx, c, "missing authorization header", constants.ERROR_UNAUTHORIZED)
			return
		}

		// Separa el usuario y la contraseña
		userPass := strings.SplitN(string(decoded), ":", 2)
		if len(userPass) != 2 {
			ins_log.Errorf(ctx, "Invalid authorization format")
			sendErrorResponse(ctx, c, "missing authorization header", constants.ERROR_UNAUTHORIZED)
			return
		}

		username := userPass[0]
		password := userPass[1]
		ins_log.Infof(ctx, "starting Authentication proccess for %s", username)

		// Aquí validamos que el usuario y la contraseña sea la correcta por ahora quemamos resultado
		userData, err := getUserData(username, ctx)
		if err != nil {
			ins_log.Errorf(ctx, "error getting user data to the user %s", username)
			ins_log.Errorf(ctx, "Invalid username or password")
			sendErrorResponse(ctx, c, "Invalid username or password", constants.ERROR_INVALID_USERNAME_OR_PASSWORD)
			return
		}

		//formateamos la password
		formatedPassword := formatPassword(password)
		ins_log.Tracef(ctx, "formatted password: %s", formatedPassword)

		if username != userData.Username || formatedPassword != userData.Password {
			ins_log.Errorf(ctx, "Invalid username or password")
			sendErrorResponse(ctx, c, "Invalid username or password", constants.ERROR_INVALID_USERNAME_OR_PASSWORD)
			return
		}

		ins_log.Infof(ctx, "the user:%s is already authenticated", username)

		//guardamos el nombre de usarrio en el contexto y tambien el dominio que luego utilizaremos
		c.Set("username", userData.Username)
		c.Set("dominio", userData.UserDomain)

		// Si la validación es exitosa, continúa con el próximo handler
		c.Next()
	}
}

// CALL THE CLIENT TO THE DATABASE
func getUserData(username string, ctx context.Context) (response.UserResponse, error) {

	//creamos el request que tendra el username para enviar a la base
	request := request.NewGetUserRequest(username)
	userData, err := client.CallToGetUserData(*request, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error getting user data in the database: %v", err)
		return userData, err
	}
	return userData, nil
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
