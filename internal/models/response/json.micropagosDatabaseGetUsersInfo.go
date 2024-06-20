package response

type GetUsersInfoResponse struct {
	Result  uint8      `json:"result"`
	Message string     `json:"message"`
	Users   []UserInfo `json:"users"`
}

type UserInfo struct {
	UserId          string `json:"user_id"`
	Username        string `json:"user_username"`
	UserPassword    string `json:"user_password"`
	UserLastLogin   string `json:"user_last_login"`
	DomainName      string `json:"domain_name"`
	UserUpdateLogin string `json:"user_update_login"`
}
