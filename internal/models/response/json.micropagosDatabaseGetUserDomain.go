package response

type UserDomainResponse struct {
	Result     uint   `json:"result"`
	Message    string `json:"message"`
	UserDomain string `json:"domain"`
}
