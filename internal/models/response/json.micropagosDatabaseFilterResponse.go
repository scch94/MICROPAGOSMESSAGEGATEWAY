package response

type FilterResponse struct {
	IsFilter bool   `json:"result"`
	Message  string `json:"message"`
}
