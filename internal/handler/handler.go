package handler

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/client"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/helper"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/request"
	"github.com/scch94/ins_log"
)

type Handler struct {
	utfi string
}

func NewHandler() *Handler {
	handler := &Handler{}
	handler.utfi = ins_log.GenerateUTFI()
	return handler
}

func SendmassiveMessage(r *request.SendMessageRequest, username string, ctx context.Context) chan helper.Result {

	//creamos la structura que guardara los resultados de cada proceso-> slices
	ins_log.Info(ctx, "starting to send massive message")

	responseErrors := make(chan helper.Result, len(r.Body.SendMassiveMessages.MobileMessageDto))
	var wg sync.WaitGroup
	var i int

	// Ciclo for que recorrera los mensajes a enviar
	for identifier, p := range r.Body.SendMassiveMessages.MobileMessageDto {
		wg.Add(1)
		//additional utfi para identiicar los distintos proceso
		utfi := ins_log.GenerateUTFI()
		ins_log.Infof(ctx, "this is the individual identifiar of this petition %v", utfi)
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
			ShortNumber:    r.Body.SendMassiveMessages.ShortNumber,
			Result:         "",
			StartPetition:  time.Now(),
		}
		go func(validate *helper.ToValidate, i int) {
			defer wg.Done()
			//creamos la structura de guardar el resultado
			result := helper.Result{}
			ins_log.Debugf(ctx, "[%v], starting to validate the petition ", utfi)

			// Validamos el mensaje
			validationResult := SendSingleMessage(&validationStruct, utfi, ctx)
			ins_log.Tracef(ctx, "[%v], this is the data of the validation result %+v", utfi, validationResult)
			module, err := validationResult.SearchValidationResultError()
			if err != nil {
				ins_log.Debugf(ctx, "[%v], error when we try to send the message time to insert the message into the database err: %v", utfi, err)

			} else {
				ins_log.Tracef(ctx, "[%v], No error time to insert the message into the database ", utfi)
			}

			//llenamos el result con el valor que insetaremos en la base de datos en el status
			validationStruct.Result = module

			//ahora vamos a insertar el mensaje a la base de datos !
			insertResult := client.CallToInsertMessageDB(&validationStruct, utfi, ctx)
			if insertResult.Id == "" {
				ins_log.Errorf(ctx, "[%v], error calling database err: %v", utfi, err)
			} else {
				ins_log.Tracef(ctx, "[%v], transaction with id %v inserted correctly!", utfi, insertResult.Id)
			}
			result.ValidationResult = validationResult
			result.InsertResult = insertResult

			ins_log.Tracef(ctx, "[%v] Ended", utfi)
			responseErrors <- result

		}(&validationStruct, i)
	}

	go func() {
		wg.Wait()
		close(responseErrors) // Cerrar el canal una vez todas las goroutines han terminado
	}()
	// Cerrar canal de errores una vez completadas las goroutines

	ins_log.Tracef(ctx, "this is the number of procced message %d", len(r.Body.SendMassiveMessages.MobileMessageDto))
	return responseErrors
}

func SendMethod(r *request.SendMessageRequest, username string, ctx context.Context) helper.Result {

	//guardamos el utfi que vamos a utilizar para seguir la sesion
	utfi := ins_log.GenerateUTFI()

	//creamos la structura de guardar errores
	result := helper.Result{}
	ins_log.Info(ctx, "method: send indivuald message, starting to validate the petition ")

	// Creamos el helper de validacion
	validationStruct := helper.NewPetition(username, r.Body.Send.Mobile, r.Body.Send.Message, r.Body.Send.UseOriginName)

	// Validamos el mensaje
	validationResult := SendSingleMessage(validationStruct, utfi, ctx)
	ins_log.Tracef(ctx, "[%v], this is the data of the validation result %+v", utfi, validationResult)
	status, err := validationResult.SearchValidationResultError()
	if err != nil {
		ins_log.Debugf(ctx, "[%v], error when we try to procees the message time to insert the message into the database with err: %v", utfi, err)

	} else {
		ins_log.Tracef(ctx, "[%v], No error time to insert the message into the database", utfi)
	}

	//llenamos el result con el valor que insetaremos en la base de datos en el status
	validationStruct.Result = status

	//ahora vamos a insertar el mensaje a la base de datos !
	insertResult := client.CallToInsertMessageDB(validationStruct, utfi, ctx)
	if insertResult.Id == "" {
		ins_log.Errorf(ctx, "[%v], error calling database err: %v", utfi, err)
	} else {
		ins_log.Tracef(ctx, "[%v], transaction with id %v inserted correctly!", utfi, insertResult.Id)
	}
	result.ValidationResult = validationResult
	result.InsertResult = insertResult
	ins_log.Tracef(ctx, "[%v] Ended", utfi)
	return result
}
func SendSingleMessage(validate *helper.ToValidate, utfi string, ctx context.Context) helper.ValidationResult {

	//creamos la struct devalidacion
	validationResult := helper.ValidationResult{}

	ins_log.Debugf(ctx, "[%v] this is the data that we are going to validate : %s", utfi, validate.ToString())

	//1era validacion mobileregex
	err := validate.ValidateMobileRegex(utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "[%v], error in the function validateMobileRegex()", utfi)
		validationResult.PassedValidation = false
		validationResult.ValidationMessage = err.Error()
		return validationResult
	}
	ins_log.Debugf(ctx, "[%v], the mobile number pass the regex expression and the formatted number is %v", utfi, validate.Mobile)

	//2do validamos el largo del mensaje y vemos si usa el massive message o si tiene un mensaje definido
	err = validate.ValidateMessageLength(utfi, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "[%v], error in ValidateMessageLengt(): ", utfi)
		validationResult.PassedValidation = false
		validationResult.ValidationMessage = err.Error()
		return validationResult
	}
	ins_log.Debugf(ctx, "[%v], the message pass the validateMessageLength this is the final message: %v", utfi, validate.Message)
	validationResult.PassedValidation = true
	validationResult.ValidationMessage = ""

	//3ero validateShortNumber si el shortnumber esta en la peticion o usamos el dominio de la confgi!
	userDomainResult := GetShortNumber(validate, utfi, ctx)
	if userDomainResult.UserDomainError != nil {

		ins_log.Errorf(ctx, "[%v], error when we try to getshortnumber(): ", utfi)
		validationResult.UserDomainResult = userDomainResult
		return validationResult
	}
	validationResult.UserDomainResult = userDomainResult
	ins_log.Infof(ctx, "[%v] this is the originNumber %v", utfi, validate.ShortNumber)

	//4to Obtenemos el telcoName llamando a portabilidad, internamente llamara tambien a
	portabilidadResult := client.CallPortabilidad(validate, utfi, ctx)
	if !portabilidadResult.PassedPortabilidad {
		ins_log.Errorf(ctx, "[%v], error in callportabilidad()", utfi)
		validationResult.PortabilidadResult = portabilidadResult

		return validationResult
	}
	validationResult.PortabilidadResult = portabilidadResult
	ins_log.Infof(ctx, "[%v], TELCONAME: %s", utfi, validate.Telco)

	//5TO FILTERED VEMOS SI EL ORIGIN Y DESTINO ESTAN FILTRADOS!
	ins_log.Tracef(ctx, "[%v], starting to validate if the the origin number and destiny are filters", utfi)
	filterResult := client.CallToFiltersDB(validate, utfi, ctx)
	validationResult.FilterResult = filterResult
	if filterResult.Error != nil {
		ins_log.Errorf(ctx, "[%v], error in callFilterdb():", utfi)
		telcoGatewayResullt := helper.SmsgatewayResult{
			PassedSmsgateway:      true,
			SmsgatewayResult:      "0",
			SmsgatewayDescription: "passed",
		}
		validationResult.SmsgatewayResult = telcoGatewayResullt
		return validationResult
	}
	if filterResult.IsFilter {
		ins_log.Tracef(ctx, "[%v], the origin number and the destinity are filters and the reason is: %v", utfi, filterResult.FilterMessage)
		telcoGatewayResullt := helper.SmsgatewayResult{
			PassedSmsgateway:      true,
			SmsgatewayResult:      "0",
			SmsgatewayDescription: "passed",
		}
		validationResult.SmsgatewayResult = telcoGatewayResullt
		return validationResult
	}

	//6to sendafter and sendbefore esta no devuelve error si encuentra algun problema los valores de sendafter o senfbefore estaran vacios
	ins_log.Tracef(ctx, "[%v], starting to check send after and send before data", utfi)
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
		ins_log.Errorf(ctx, "[%v], error in CallTelcoGateway()", utfi)
		validationResult.SmsgatewayResult = telcoGatewayResullt
		return validationResult
	}

	//guardamos la hora en la que enviamos el mensaje !
	validate.SendMessageTime = time.Now()
	ins_log.Infof(ctx, "[%v], TelcoGatewayResponse: %s", utfi, telcoGatewayResullt.SmsgatewayResult)

	//llenamos el validattionresult con los datos de smsgatewayresult t
	validationResult.SmsgatewayResult = telcoGatewayResullt

	return validationResult

}

func GetShortNumber(validate *helper.ToValidate, utfi string, ctx context.Context) helper.UserDomainResult {

	//llenamos el herper que sera utilizado para el domainresult
	userDomainResult := helper.UserDomainResult{}
	ins_log.Infof(ctx, "[%v], checking short number", utfi)
	if validate.ShortNumber != "" {

		//si el shortnumber viene en la peticion usamos ese shortnumber
		ins_log.Tracef(ctx, "[%v], using short number of the petition", utfi)
		return helper.UserDomainResult{
			UserDomainResult:  "ok",
			UserDomainError:   nil,
			UserDomainMessage: "using short of the petition",
		}

	} else {

		//si no usamos el que esta en la cofig
		ins_log.Tracef(ctx, "[%v], using short number of the config", utfi)

		//ahora con el dominio vamos a buscar el shortnumberen la config
		shortNumber := config.Config.SearchShortNumber(dominioString)
		if shortNumber == "" {
			ins_log.Tracef(ctx, "[%v], dindt find a short number for the domain %s please check the configuration", dominioString)
			userDomainResult.UserDomainError = errors.New("dindt find a short number for the domain " + dominioString)
			userDomainResult.UserDomainMessage = "dindt find a short number for the domain"
			return userDomainResult
		}
		userDomainResult = helper.UserDomainResult{
			UserDomainResult:  "ok",
			UserDomainError:   nil,
			UserDomainMessage: "using short number of the config",
		}
		validate.ShortNumber = shortNumber
		return userDomainResult
	}

}
