package response

import "encoding/json"

type InsertMessageResponse struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
	Id      int    `json:"Id"`
}

func (payload *InsertMessageResponse) ToJSON() (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
