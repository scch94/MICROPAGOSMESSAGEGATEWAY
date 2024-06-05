package helper

import (
	"regexp"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
)

var Mask []response.Mask

func GetMask(shortCode string, message string) string {

	modifiedMessage := ""
	for _, m := range Mask {
		if shortCode == m.ShortNumber && m.Direction == "OUT" {
			maskPattern := regexp.MustCompile(m.MaskPattern)
			modifiedMessage = maskPattern.ReplaceAllString(message, "X")
			return modifiedMessage
		}
	}
	modifiedMessage = message
	return modifiedMessage
}
