package handler

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

// import (
// 	"context"
// 	"encoding/xml"
// 	"errors"
// 	"io"
// 	"time"

// 	"net/http"
// 	"sync"

// 	"github.com/gin-gonic/gin"
// 	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/client"
// 	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
// 	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
// 	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
// 	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
// 	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
// 	"github.com/scch94/ins_log"
// )

func (h *Handler) LegacySendMessageHandler(c *gin.Context) {

	//traemos el contexto y le setiamos el contexto actual
	ctx := c.Request.Context()
	ctx = ins_log.SetPackageNameInContext(ctx, "handler")

	//traemos el usuarios desde el contexto no usamos el error por que si llegamos hasta aca el username debe existir
	username, _ := c.Get("username")
	dominio, _ := c.Get("dominio")

	usernameString := username.(string)
	dominioString = dominio.(string)
	ins_log.Debugf(ctx, "starting to legacysendMessage method to the user %s", username)
	ins_log.Debug(ctx, "starting to get the xml")

	//leemos el cuerpo de la peticion! y manejamos su error
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ins_log.Errorf(ctx, "Error reading the body of the petition: %s", err)
		code := "env:Server"
		message := "Error reading the body of the petition"
		xmlresponse := response.GenerateFaultXML(code, message)
		c.Data(http.StatusBadRequest, "application/xml", []byte(xmlresponse))
		return
	}

	//creamos la variable que sera del tipo struct para guardar los datos de la peticion
	var sendMessageRequest request.SendMessageRequest

	//parceamos el xml y lo guadamos en envelope
	err = xml.Unmarshal(body, &sendMessageRequest)
	if err != nil {
		ins_log.Errorf(ctx, "error unmarshalling the body of the petition err: %s", err)
		code := "env:Server"
		message := "error unmarshalling the body of the petition"
		xmlresponse := response.GenerateFaultXML(code, message)
		c.Data(http.StatusBadRequest, "application/xml", []byte(xmlresponse))
		return
	}

	//miramos que metodo es el de la solicitud si el massive message o el message legacy
	method, err := sendMessageRequest.DetermineRequestType()
	if err != nil {
		ins_log.Errorf(ctx, "error method not allowed: %s", err)
		code := "env:Server"
		message := "error method not allowed"
		xmlresponse := response.GenerateFaultXML(code, message)
		c.Data(http.StatusBadRequest, "application/xml", []byte(xmlresponse))
		return
	}

	//dependiendo el metodo comenzamos a operar.
	if method == constants.SEND {

		//vamos al send method que nos devolvera una structura con el resultado de cada una de las operaciones del envio de mensajes
		result := SendMethod(&sendMessageRequest, usernameString, ctx)

		//como la funcion createdinsertresults pide un slice hacemos un slice de una unica posision
		results := []helper.Result{result}

		//esto nos devolvera el xml para la respuesta
		xmlBody, err := createdLegacyResult(results, method, ctx)
		if err != nil {
			c.Data(http.StatusInternalServerError, "application/xml", []byte(xmlBody))
			return
		}

		c.Data(http.StatusOK, "application/xml", []byte(xmlBody))
		return
	} else {
		err = checkMobilereg(sendMessageRequest.Body.SendMassiveMessages.MobileMessageDto, ctx)
		if err != nil {
			ins_log.Errorf(ctx, "any of the movils passed the regex expression: %v", err)
			code := "env:Server"
			message := "error in the data of the petition"
			xmlresponse := response.GenerateFaultXML(code, message)
			c.Data(http.StatusInternalServerError, "application/xml", []byte(xmlresponse))
			return
		}
		//vamos al send massive method el cual nos devolvera un chan
		results := SendmassiveMessage(&sendMessageRequest, usernameString, ctx)

		//como el createdinsertresults recibe un slice de results volvemos el chan en un slice
		var resultsSlice []helper.Result
		for result := range results {
			resultsSlice = append(resultsSlice, result)
		}

		//esto nos devolvera el xml que usaremos
		xmlBody, _ := createdLegacyResult(resultsSlice, method, ctx)

		c.Data(http.StatusOK, "application/xml", []byte(xmlBody))
	}
}

func createdLegacyResult(results []helper.Result, method string, ctx context.Context) (string, error) {
	var xmlResponse string
	var err error
	if method == constants.SEND {
		for _, result := range results {
			_, err = result.SerchResultError()
			if err != nil {
				codError := "env:Server"
				message := err.Error()
				xmlResponse = response.GenerateFaultXML(codError, message)
				return xmlResponse, err
			}
		}
		xmlResponse, err = response.GenerateLegacyXML(false)
		if err != nil {
			ins_log.Errorf(ctx, "Error generating XML: %v", err)
			return "", err
		}
		return xmlResponse, nil
	} else {
		xmlResponse, err = response.GenerateLegacyXML(true)
		if err != nil {
			ins_log.Errorf(ctx, "Error generating XML: %v", err)
			return "", err
		}
	}

	return xmlResponse, nil

}

// esta funcion se encargara de ver si alguno de los numeros de la peticion pasa la expresion regular si alguno la pasa se continua con el procesos si no devuelve un error
func checkMobilereg(mobileMessages []request.MobileMessageDto, ctx context.Context) error {
	var count int
	regex, err := regexp.Compile(config.Config.MobileRegex)
	if err != nil {
		ins_log.Errorf(ctx, "error to compilate the regex expression function regexp.Compile(): , err: %v", err)
		return fmt.Errorf("error when we try to compilate the regex check the config %w", err)
	}
	for _, momobileMessage := range mobileMessages {
		if regex.MatchString(momobileMessage.Mobile) {
			ins_log.Tracef(ctx, "mobile %v match in the regex expression", momobileMessage.Mobile)
			count++
		}
	}
	if count != 0 {
		return nil
	} else {
		return fmt.Errorf("any mobile match with the regex expression")
	}
}
