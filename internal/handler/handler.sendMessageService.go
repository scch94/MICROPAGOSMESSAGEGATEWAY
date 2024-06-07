package handler

import (
	"encoding/xml"
	"io"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

var dominioString string

func (h *Handler) SendMessageHandler(c *gin.Context) {

	//traemos el contexto y le setiamos el contexto actual
	ctx := c.Request.Context()
	ctx = ins_log.SetPackageNameInContext(ctx, "handler")

	//traemos el usuarios desde el contexto no usamos el error por que si llegamos hasta aca el username debe existir
	username, _ := c.Get("username")
	dominio, _ := c.Get("dominio")

	usernameString := username.(string)
	dominioString = dominio.(string)
	ins_log.Infof(ctx, "starting to sendMessage method to the user %s", username)
	ins_log.Debug(ctx, "starting to get the xml")

	//leemos el cuerpo de la peticion! y manejamos su error
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ins_log.Errorf(ctx, "Error reading the body of the petition: %s", err)
		result := response.Result{
			DatabaseID: "",
			Message:    err.Error(),
			Code:       constants.ERROR_READING_THE_BODY,
		}
		results := []response.Result{result}
		xmlresponse := response.GenerateXML(results)
		c.Data(http.StatusBadRequest, "application/xml", []byte(xmlresponse))
		return
	}
	//creamos la variable que sera del tipo struct para guardar los datos de la peticion
	var SendMessageRequest request.SendMessageRequest

	//parceamos el xml y lo guadamos en envelope
	err = xml.Unmarshal(body, &SendMessageRequest)
	if err != nil {
		ins_log.Errorf(ctx, "error unmarshalling the body of the petition err: %s", err)
		result := response.Result{
			DatabaseID: "",
			Message:    err.Error(),
			Code:       constants.ERROR_UNMARSHALL_THE_BODY,
		}
		results := []response.Result{result}
		xmlresponse := response.GenerateXML(results)
		c.Data(http.StatusBadRequest, "application/xml", []byte(xmlresponse))
		return
	}

	//miramos que metodo es el de la solicitud si el massive message o el message legacy
	method, err := SendMessageRequest.DetermineRequestType()
	if err != nil {
		result := response.Result{
			DatabaseID: "",
			Message:    err.Error(),
			Code:       constants.ERROR_CHECKING_REQUEST_TYPE,
		}
		results := []response.Result{result}
		xmlresponse := response.GenerateXML(results)
		c.Data(http.StatusBadRequest, "application/xml", []byte(xmlresponse))
		return
	}

	//dependiendo el metodo comenzamos a operar.
	if method == constants.SEND {

		//vamos al send method que nos devolvera una structura con el resultado de cada una de las operaciones del envio de mensajes
		result := SendMethod(&SendMessageRequest, usernameString, ctx)

		//como la funcion createdinsertresults pide un slice hacemos un slice de una unica posision
		results := []helper.Result{result}

		//esto nos devolvera un arreglo con los datos que tendra el body
		xmlBody := createdInsertResults(results)

		//aquii generamos un string de forma de xml para enviarlo
		response := response.GenerateXML(xmlBody)
		c.Data(http.StatusOK, "application/xml", []byte(response))
		return
	} else {

		//vamos al send massive method el cual nos devolvera un chan
		results := SendmassiveMessage(&SendMessageRequest, usernameString, ctx)

		//como el createdinsertresults recibe un slice de results volvemos el chan en un slice
		var resultsSlice []helper.Result
		for result := range results {
			resultsSlice = append(resultsSlice, result)
		}

		//esto nos devolvera un arreglo con los datos que tendra el body
		xmlBody := createdInsertResults(resultsSlice)
		//aquii generamos un string de forma de xml para enviarlo
		response := response.GenerateXML(xmlBody)
		c.Data(http.StatusOK, "application/xml", []byte(response))
	}
}

// utils

func createdInsertResults(results []helper.Result) []response.Result {
	var responses []response.Result
	for _, result := range results {
		var response response.Result
		code, err := result.SerchResultError()
		if err != nil {
			response.DatabaseID = result.InsertResult.Id
			response.Message = err.Error()
			response.Code = code
		} else {
			if code == constants.OK {
				response.DatabaseID = result.InsertResult.Id
				response.Message = "meesagge inserted correctly"
				response.Code = code
			} else if code == constants.ERROR_USER_IS_FILTERED {
				response.DatabaseID = result.InsertResult.Id
				response.Message = "meesagge inserted incorrectly with status FILTERED and we are not allowed to send this message"
				response.Code = code
			} else {
				response.DatabaseID = result.InsertResult.Id
				response.Message = "meesagge inserted correctly we are going to send the message later"
				response.Code = code
			}

		}
		responses = append(responses, response)
	}
	return responses
}
