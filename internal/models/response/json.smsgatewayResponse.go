package response

import "encoding/json"

type SmsGatewayResponse struct {
	Status      string `json:"status"`
	ForwardRef  string `json:"forwardRef"`
	Description string `json:"Description"`
}

func (payload *SmsGatewayResponse) ToJSON() (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
