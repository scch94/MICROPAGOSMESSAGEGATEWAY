package constants

const (
	PATH               = "ws/SendMessageService"
	MESSAGE_REGEX      = "[áéíóúÁÉÍÓÚñÑ]"
	MAX_MESSAGE_LENGTH = 500
	GATEWAY_URL        = "http://localhost:1010"

	//constantest usadas en la validacion
	MASSIVE_MESSAGE = "MassiveMessage"
	SEND            = "Send"
	CHECKTIMEFORMAT = "15:04"
	OUTPUTFORMAT    = "15:04:05"

	//constants usadas en el llamado para insetar en la base de datos
	INSERT_TYPE                     = "OutboundSMS"
	INSERT_MOBILEC_OUNTRY_ISO_CODE  = "UY"
	INSERT_SESSION_ACTION           = "NONE"
	INSERT_SESSION_APPLICATION_NAME = "micropagos 2.0"

	//constants usadas para identificar errors dentro del sistema
	ERROR_READING_THE_BODY                     = "1"
	ERROR_UNMARSHALL_THE_BODY                  = "2"
	ERROR_CHECKING_REQUEST_TYPE                = "3"
	ERROR_VALIDATION_DATA_OF_THE_PETITION      = "4"
	ERROR_GETTING_USER_DOMAIN                  = "5"
	ERROR_GETTING_PORTABILITY                  = "6"
	ERROR_GETTING_FILTER                       = "7"
	ERROR_USER_IS_FILTERED                     = "8"
	ERROR_CHECKING_WHEN_SEND_MESSAGE           = "9"
	ERROR_IN_SMS_GATEWAY                       = "10"
	ERROR_INSERTING_MESSAGE_RESULT_IN_DATABASE = "11"
	ERROR_UNAUTORIES                           = "12"
	OK                                         = "0"

	//RESULTADOS PARA INSERTAR EN LA BASE DE DATOS
	RESULT_ERROR                 = "ERROR"
	RESULT_ERROR_INVOKING_SENDER = "ERROR_INVOKING_SENDER"
	RESULT_FILTERED              = "FILTERED"
	RESULT_PENDING               = "PENDING"
	RESULT_SENT                  = "SENT"
)
