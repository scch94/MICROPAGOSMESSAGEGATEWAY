package helper

import (
	"errors"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
)

type Result struct {
	ValidationResult
	InsertResult
}

type ValidationResult struct {
	PassedValidation  bool
	ValidationMessage string
	UserDomainResult
	PortabilidadResult
	FilterResult
	WhenSendResult
	SmsgatewayResult
}

type UserDomainResult struct {
	UserDomainResult  string
	UserDomainMessage string
	UserDomainError   error
}

type PortabilidadResult struct {
	PassedPortabilidad  bool
	PortabilidadMessage string
}

type FilterResult struct {
	IsFilter      bool
	FilterMessage string
	Error         error
}
type WhenSendResult struct {
	SendNow     bool
	SendMessage string
}

type SmsgatewayResult struct {
	PassedSmsgateway      bool
	SmsgatewayResult      string
	SmsgatewayDescription string
}

type InsertResult struct {
	Result  string
	Message string
	Id      string
}

// este metodo devolvera si existio un error o no durante los resultados si encuentra error devuelve el error y en el nombre la funcion de validacion que se rompio importante validar el orden de los errores
func (v *ValidationResult) SearchValidationResultError() (string, error) {
	if !v.PassedValidation {
		return constants.RESULT_ERROR, errors.New(v.ValidationMessage)
	}
	if v.UserDomainResult.UserDomainError != nil {
		return constants.RESULT_ERROR, errors.New(v.UserDomainResult.UserDomainMessage)
	}
	if !v.PortabilidadResult.PassedPortabilidad {
		return constants.RESULT_ERROR_INVOKING_SENDER, errors.New(v.PortabilidadResult.PortabilidadMessage)
	}
	if v.FilterResult.Error != nil {
		return constants.RESULT_ERROR_INVOKING_SENDER, v.FilterResult.Error
	}
	if v.FilterResult.IsFilter {
		return constants.ERROR_USER_IS_FILTERED, nil
	}
	if !v.WhenSendResult.SendNow {
		return constants.RESULT_PENDING, nil
	}
	if !v.SmsgatewayResult.PassedSmsgateway {
		return constants.RESULT_ERROR_INVOKING_SENDER, errors.New("error sending SMSgateway")
	}

	return constants.RESULT_SENT, nil
}

func (r *Result) SerchResultError() (string, error) {
	if !r.ValidationResult.PassedValidation {
		return constants.ERROR_VALIDATION_DATA_OF_THE_PETITION, errors.New("error when we try to validate the petition")
	}
	if r.UserDomainResult.UserDomainError != nil {
		return constants.ERROR_GETTING_USER_DOMAIN, errors.New("error when we try to get the origin number")
	}
	if !r.ValidationResult.PortabilidadResult.PassedPortabilidad {
		return constants.ERROR_GETTING_PORTABILITY, errors.New("error when we go to the portability module")
	}
	if r.ValidationResult.FilterResult.Error != nil {
		return constants.ERROR_GETTING_FILTER, errors.New("error when we try to see if destini is filtered")
	}
	if r.ValidationResult.FilterResult.IsFilter {
		return constants.ERROR_USER_IS_FILTERED, nil
	}
	if !r.ValidationResult.WhenSendResult.SendNow {
		return constants.ERROR_CHECKING_WHEN_SEND_MESSAGE, nil
	}
	if !r.ValidationResult.SmsgatewayResult.PassedSmsgateway {
		return constants.ERROR_IN_SMS_GATEWAY, errors.New("error when we go to the smsgateway module")
	}
	if r.InsertResult.Id == "" {
		return constants.ERROR_INSERTING_MESSAGE_RESULT_IN_DATABASE, errors.New("error when we try to insert the message in the database")
	}

	return constants.OK, nil
}
