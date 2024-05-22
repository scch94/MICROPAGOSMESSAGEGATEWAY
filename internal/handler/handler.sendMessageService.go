package handler

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"time"

	"net/http"
	"sync"

	//"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/client"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

func (h *Handler) SendMessageService(c *gin.Context) {

	//traemos el contexto y le setiamos el contexto actual
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, constants.PACKAGE_NAME_KEY, "handler")

	//traemos el usuarios desde el contexto no usamos el error por que si llegamos hasta aca el username debe existir
	username, _ := c.Get("username")
	usernameString := username.(string)
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
	xmlString := string(body)
	ins_log.Tracef(ctx, "this is the body of the petition %v", xmlString)

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
		result := sendMethod(&SendMessageRequest, usernameString, ctx)

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
		results := sendmassiveMessage(&SendMessageRequest, usernameString, ctx)

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

func sendmassiveMessage(r *request.SendMessageRequest, username string, ctx context.Context) chan helper.Result {

	//creamos la structura que guardara los resultados de cada proceso-> slices
	ins_log.Info(ctx, "starting to send massive message")
	var wg sync.WaitGroup
	responseErrors := make(chan helper.Result, len(r.Body.SendMassiveMessages.MobileMessageDto))
	var i int

	// Ciclo for que recorrera los mensajes a enviar
	for identifier, p := range r.Body.SendMassiveMessages.MobileMessageDto {

		//additional utfi para identiicar los distintos proceso
		utfi := ins_log.GenerateUTFI()
		ins_log.Infof(ctx, "this is the identifier of this petition %v", utfi)
		i = identifier
		// Creamos el helper de validación
		validationStruct := helper.ToValidate{
			Username:       username,
			Mobile:         p.Mobile,
			Message:        p.Message,
			UseOriginName:  r.Body.SendMassiveMessages.UseOriginName,
			MassiveMessage: r.Body.SendMassiveMessages.MassiveMessage,
			SendAfter:      r.Body.SendMassiveMessages.SendAfter,
			SendBefore:     r.Body.SendMassiveMessages.SendBefore,
			Priority:       r.Body.SendMassiveMessages.Priority,
			Result:         "",
			StartPetition:  time.Now(),
		}
		wg.Add(1)
		go func(validate *helper.ToValidate, i int) {

			//creamos la structura de guardar el resultado
			result := helper.Result{}
			ins_log.Infof(ctx, "PETITION[%v], starting to validate the petition #%v", utfi)

			// Validamos el mensaje
			validationResult := sendSingleMessage(&validationStruct, utfi, ctx)
			ins_log.Tracef(ctx, "PETITION[%v], this is the data of the validation result %+v", utfi, validationResult)
			module, err := validationResult.SearchValidationResultError()
			if err != nil {
				ins_log.Debugf(ctx, "PETITION[%v], error when we try to send the message time to insert the message into the database err: %v", utfi, err)

			} else {
				ins_log.Tracef(ctx, "PETITION[%v], No error time to insert the message into the database ", utfi)
			}

			//llenamos el result con el valor que insetaremos en la base de datos en el status
			validationStruct.Result = module

			//ahora vamos a insertar el mensaje a la base de datos !
			insertResult := client.CallToInsertMessageDB(&validationStruct, utfi, ctx)
			if insertResult.Id == "" {
				ins_log.Errorf(ctx, "PETITION[%v], error calling database err: %v", utfi, err)
			} else {
				ins_log.Tracef(ctx, "PETITION[%v], transaction with id %v inserted correctly!", utfi, insertResult.Id)
			}
			result.ValidationResult = validationResult
			result.InsertResult = insertResult

			ins_log.Tracef(ctx, "PETITION[%v] Ended", utfi)
			responseErrors <- result

			wg.Done()
		}(&validationStruct, i)
	}
	wg.Wait()
	close(responseErrors) // Cerrar canal de errores una vez completadas las goroutines

	ins_log.Infof(ctx, "this is the number of procced message %d", len(r.Body.SendMassiveMessages.MobileMessageDto))
	return responseErrors
}

func sendMethod(r *request.SendMessageRequest, username string, ctx context.Context) helper.Result {

	//guardamos el utfi que vamos a utilizar para seguir la sesion
	utfi := ins_log.GenerateUTFI()

	//creamos la structura de guardar errores
	result := helper.Result{}
	ins_log.Info(ctx, "method: send indivuald message, starting to validate the petition ")

	// Creamos el helper de validacion
	validationStruct := helper.NewPetition(username, r.Body.Send.Mobile, r.Body.Send.Message, r.Body.Send.UseOriginName)

	// Validamos el mensaje
	validationResult := sendSingleMessage(validationStruct, utfi, ctx)
	ins_log.Tracef(ctx, "PETITION[%v], this is the data of the validation result %+v", utfi, validationResult)
	status, err := validationResult.SearchValidationResultError()
	if err != nil {
		ins_log.Debugf(ctx, "PETITION[%v], error when we try to procees the message time to insert the message into the database with err: %v", utfi, err)

	} else {
		ins_log.Tracef(ctx, "PETITION[%v], No error time to insert the message into the database", utfi)
	}

	//llenamos el result con el valor que insetaremos en la base de datos en el status
	validationStruct.Result = status

	//ahora vamos a insertar el mensaje a la base de datos !
	insertResult := client.CallToInsertMessageDB(validationStruct, utfi, ctx)
	if insertResult.Id == "" {
		ins_log.Errorf(ctx, "PETITION[%v], error calling database err: %v", utfi, err)
	} else {
		ins_log.Tracef(ctx, "PETITION[%v], transaction with id %v inserted correctly!", utfi, insertResult.Id)
	}
	result.ValidationResult = validationResult
	result.InsertResult = insertResult
	ins_log.Tracef(ctx, "PETITION[%v] Ended", utfi)
	return result
}
func sendSingleMessage(validate *helper.ToValidate, utfi string, ctx context.Context) helper.ValidationResult {

	//creamos la struct devalidacion
	validationResult := helper.ValidationResult{}

	ins_log.Debugf(ctx, "PETITION[%v] this is the data that we are going to validate : %s", utfi, validate.ToString())

	//1era validacion mobileregex
	err := validate.ValidateMobileRegex(utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error in the function validateMobileRegex()", utfi)
		validationResult.PassedValidation = false
		validationResult.ValidationMessage = err.Error()
		return validationResult
	}
	ins_log.Debugf(ctx, "PETITION[%v], the mobile number pass the regex expression and the formatted number is %v", utfi, validate.Mobile)

	//2do validamos el largo del mensaje y vemos si usa el massive message o si tiene un mensaje definido
	err = validate.ValidateMessageLength(utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error in ValidateMessageLengt(): ", utfi)
		validationResult.PassedValidation = false
		validationResult.ValidationMessage = err.Error()
		return validationResult
	}
	ins_log.Debugf(ctx, "PETITION[%v], the message pass the validateMessageLength this is the final message: %v", utfi, validate.Message)
	validationResult.PassedValidation = true
	validationResult.ValidationMessage = ""

	//3ero validateShortNumber si el shortnumber esta en la peticion no vamos a la base y si no esta vamos a la base !
	userDomainResult := GetShortNumber(validate, utfi, ctx)
	if userDomainResult.UserDomainError != nil {

		ins_log.Errorf(ctx, "PETITION[%v], error when we try to getshortnumber(): ", utfi)
		validationResult.UserDomainResult = userDomainResult
		return validationResult
	}
	validationResult.UserDomainResult = userDomainResult
	ins_log.Infof(ctx, "PETITION[%v] this is the originNumber %v", utfi, validate.ShortNumber)

	//4to Obtenemos el telcoName llamando a portabilidad, internamente llamara tambien a
	portabilidadResult := client.CallPortabilidad(validate, utfi, ctx)
	if !portabilidadResult.PassedPortabilidad {
		ins_log.Errorf(ctx, "PETITION[%v], error in callportabilidad()", utfi)
		validationResult.PortabilidadResult = portabilidadResult

		return validationResult
	}
	validationResult.PortabilidadResult = portabilidadResult
	ins_log.Infof(ctx, "PETITION[%v], TELCONAME: %s", utfi, validate.Telco)

	//5TO FILTERED VEMOS SI EL ORIGIN Y DESTINO ESTAN FILTRADOS!
	ins_log.Tracef(ctx, "PETITION[%v], starting to validate if the the origin number and destiny are filters", utfi)
	filterResult := client.CallToFiltersDB(validate, utfi, ctx)
	validationResult.FilterResult = filterResult
	if filterResult.Error != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error in callFilterdb():", utfi)
		telcoGatewayResullt := helper.SmsgatewayResult{
			PassedSmsgateway:      true,
			SmsgatewayResult:      "0",
			SmsgatewayDescription: "passed",
		}
		validationResult.SmsgatewayResult = telcoGatewayResullt
		return validationResult
	}
	if filterResult.IsFilter {
		ins_log.Tracef(ctx, "PETITION[%v], the origin number and the destinity are filters and the reason is: %v", utfi, filterResult.FilterMessage)
		telcoGatewayResullt := helper.SmsgatewayResult{
			PassedSmsgateway:      true,
			SmsgatewayResult:      "0",
			SmsgatewayDescription: "passed",
		}
		validationResult.SmsgatewayResult = telcoGatewayResullt
		return validationResult
	}

	//6to sendafter and sendbefore esta no devuelve error si encuentra algun problema los valores de sendafter o senfbefore estaran vacios
	ins_log.Tracef(ctx, "PETITION[%v], starting to check send after and send before data", utfi)
	whenSendResult := validate.ValidateSendAfterAndSendBefore(utfi, ctx)
	if validate.SendAfter != "" || validate.SendBefore != "" {
		telcoGatewayResullt := helper.SmsgatewayResult{
			PassedSmsgateway:      true,
			SmsgatewayResult:      "0",
			SmsgatewayDescription: "passed",
		}
		validationResult.WhenSendResult = whenSendResult
		validationResult.SmsgatewayResult = telcoGatewayResullt
		return validationResult
	}
	validationResult.WhenSendResult = whenSendResult

	//6to llamamos al gateway de envio de mensajes
	telcoGatewayResullt := client.CallTelcoGateway(validate, utfi, ctx)
	if telcoGatewayResullt.SmsgatewayResult != "0" {
		ins_log.Errorf(ctx, "PETITION[%v], error in CallTelcoGateway()", utfi)
		validationResult.SmsgatewayResult = telcoGatewayResullt
		return validationResult
	}

	//guardamos la hora en la que enviamos el mensaje !
	validate.SendMessageTime = time.Now()
	ins_log.Infof(ctx, "PETITION[%v], TelcoGatewayResponse: %s", utfi, telcoGatewayResullt.SmsgatewayResult)

	//llenamos el validattionresult con los datos de smsgatewayresult t
	validationResult.SmsgatewayResult = telcoGatewayResullt

	return validationResult

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

func GetShortNumber(validate *helper.ToValidate, utfi string, ctx context.Context) helper.UserDomainResult {

	ins_log.Infof(ctx, "PETITION[%v], checking short number", utfi)
	if validate.ShortNumber != "" {

		//si el shortnumber viene en la peticion usamos ese shortnumber
		ins_log.Tracef(ctx, "PETITION[%v], using short number of the petition", utfi)
		return helper.UserDomainResult{
			UserDomainResult:  "ok",
			UserDomainError:   nil,
			UserDomainMessage: "using short of the petition",
		}

	} else {

		//si no usamos el que esta en la cofig
		ins_log.Tracef(ctx, "PETITION[%v], using short number of the config", utfi)

		//primero tenemos que ir a la base a ver que dominio tiene el usuario
		userDomainResult := client.CallToGetUserDomain(validate, utfi, ctx)
		if userDomainResult.UserDomainResult == "" || userDomainResult.UserDomainError != nil {
			ins_log.Tracef(ctx, "PETITION[%v], have an error when we try to CallToGetUserDomain() ", utfi)
			userDomainResult.UserDomainError = errors.New("have an error when we try to CallToGetUserDomain()")
			userDomainResult.UserDomainMessage = "have an error when we try to CallToGetUserDomain()"
			return userDomainResult
		}

		//ahora con el dominio vamos a buscar el shortnumberen la config
		shortNumber := config.Config.SearchShortNumber(userDomainResult.UserDomainResult)
		if shortNumber == "" {
			ins_log.Tracef(ctx, "PETITION[%v], dindt find a short number for the domain %s please check the configuration", userDomainResult.UserDomainResult)
			userDomainResult.UserDomainError = errors.New("dindt find a short number for the domain " + userDomainResult.UserDomainResult)
			userDomainResult.UserDomainMessage = "dindt find a short number for the domain"
			return userDomainResult
		}

		validate.ShortNumber = shortNumber
		return userDomainResult
	}

}
