package helper

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/config"
	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/constants"
	"github.com/scch94/ins_log"
)

type ToValidate struct {
	Username        string
	Mobile          string
	Message         string
	UseOriginName   string
	MassiveMessage  string
	SendAfter       string
	SendBefore      string
	ShortNumber     string
	Priority        string
	Telco           string
	Result          string
	SendMessageTime time.Time
	StartPetition   time.Time
}

// constructors
func NewPetition(username string, mobile string, message string, useoriginame string) *ToValidate {
	return &ToValidate{
		Username:      username,
		Mobile:        mobile,
		Message:       message,
		UseOriginName: useoriginame,
		StartPetition: time.Now(),
	}
}

func (p *ToValidate) ToString() string {
	data := ""
	if p.Mobile != "" {
		data += "Mobile: " + p.Mobile + ", "
	}
	if p.Message != "" {
		data += "Message: " + p.Message + ", "
	}
	if p.UseOriginName != "" {
		data += "UseOriginName: " + p.UseOriginName + ", "
	}
	if p.MassiveMessage != "" {
		data += "MassiveMessage: " + p.MassiveMessage + ", "
	}
	if p.SendAfter != "" {
		data += "SendAfter: " + p.SendAfter + ", "
	}
	if p.SendBefore != "" {
		data += "SendBefore: " + p.SendBefore + ", "
	}
	if p.ShortNumber != "" {
		data += "ShortNumber: " + p.ShortNumber + ", "
	}
	if p.Priority != "" {
		data += "Priority: " + p.Priority + ", "
	}

	return data
}

func (p *ToValidate) ValidateMobileRegex(utfi string, ctx context.Context) error {

	//traemos el contexto y le setiamos el contexto actual
	ctx = context.WithValue(ctx, constants.PACKAGE_NAME_KEY, "validation")

	ins_log.Infof(ctx, "PETITION[%v], starting to validate a regex expression", utfi)
	ins_log.Tracef(ctx, "PETITION[%v], regex expression: %v value to validate %v", utfi, config.Config.MobileRegex, p.Mobile)
	regex, err := regexp.Compile(config.Config.MobileRegex)
	if err != nil {
		ins_log.Errorf(ctx, "PETITION[%v], error to compilate the regex expression function regexp.Compile(): , err: %v", utfi, err)
		return fmt.Errorf("error al compilar la expression regular porfavor verificala e inteta de neuvo regexp.Compile(): %w", err)
	}
	if !regex.MatchString(p.Mobile) {
		ins_log.Errorf(ctx, "PETITION[%v], value did not match in the regex expression", utfi)
		return fmt.Errorf("value did not match in the regex expression regex.MatchString(): %w", err)
	}
	//si llegamos aca la expresion regular si valido todo.
	ins_log.Tracef(ctx, "PETITION[%v], value match with the regex expression", utfi)
	//formateamos el numero recordemos que tanto para el sms gateway como para el modulo de portabilidad el numero ira con el formato internacional
	p.Mobile = formatNumber(p.Mobile, ctx)
	return nil
}

func formatNumber(number string, ctx context.Context) string {
	ins_log.Debug(ctx, "starting to format number")
	if number[0] == '0' {
		number = number[1:]
		number = "598" + number
	}
	return number

}

func (p *ToValidate) ValidateMessageLength(utfi string, ctx context.Context) error {

	//traemos el contexto y le setiamos el contexto actual
	ctx = context.WithValue(ctx, constants.PACKAGE_NAME_KEY, "validation")

	var finalMessage string
	ins_log.Infof(ctx, "PETITION[%v], starting to validate th message length", utfi)
	ins_log.Tracef(ctx, "PETITION[%v], max length: %v message to validate %v", utfi, config.Config.MaxMessageLength, p.Message)

	if len(p.Message) == 0 && len(p.MassiveMessage) != 0 {
		ins_log.Tracef(ctx, "PETITION[%v], using massive message", utfi)
		finalMessage = p.MassiveMessage
	} else {
		ins_log.Tracef(ctx, "PETITION[%v], using normal message", utfi)
		finalMessage = p.Message
	}
	if len(finalMessage) > config.Config.MaxMessageLength {
		ins_log.Errorf(ctx, "PETITION[%v], the message length is too large the message length is %w and the max length is%w", utfi, len(finalMessage), constants.MAX_MESSAGE_LENGTH)
		return fmt.Errorf("the message length is too large the message length is %v and the max length is%v", len(finalMessage), constants.MAX_MESSAGE_LENGTH)
	}
	p.Message = finalMessage
	return nil
}

func (p *ToValidate) ValidateSendAfterAndSendBefore(utfi string, ctx context.Context) WhenSendResult {

	//traemos el contexto y le setiamos el contexto actual
	ctx = context.WithValue(ctx, constants.PACKAGE_NAME_KEY, "validation")

	whenSendResult := WhenSendResult{}
	//validaremos el send after si alguno falta el mensaje se enviara automaticamente al igual que si alguno de los dos no estan en el formato adecuado
	ins_log.Debugf(ctx, "PETITION[%v], checking send after and send before ", utfi)

	//parsing data!
	_, errSAfter := time.Parse(constants.CHECKTIMEFORMAT, p.SendAfter)
	_, errSBefore := time.Parse(constants.CHECKTIMEFORMAT, p.SendBefore)
	if errSAfter != nil || errSBefore != nil {
		whenSendResult.SendNow = true
		whenSendResult.SendMessage = "message is going to send now"
		ins_log.Tracef(ctx, "PETITION[%v] senfafter or sendbefore are empity or didnt match with the expected fomat HH:mm", utfi)
		p.SendAfter = ""
		p.SendBefore = ""
	} else {
		whenSendResult.SendNow = false
		whenSendResult.SendMessage = "message is going to send in the future"
		p.SendAfter = fmt.Sprintf(p.SendAfter + ":00")
		p.SendBefore = fmt.Sprintf(p.SendBefore + ":00")
	}
	ins_log.Debugf(ctx, "PETITION[%v]: sendBefore: %v sendAfter: %v", utfi, p.SendBefore, p.SendAfter)
	return whenSendResult
}
