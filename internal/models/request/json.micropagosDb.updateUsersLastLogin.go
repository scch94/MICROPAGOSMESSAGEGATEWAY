package request

type UsersToUpdate struct {
	Users []UserData `json:"users"`
}

type UserData struct {
	UserName  string `json:"user_name"`
	LoginTime string `json:"login_time"`
}
