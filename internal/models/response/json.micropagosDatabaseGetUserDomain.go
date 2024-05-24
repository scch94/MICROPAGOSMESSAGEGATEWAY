package response

type UserResponse struct {
	Result     uint   `json:"result"`
	Message    string `json:"message"`
	UserDomain string `json:"domain"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}
