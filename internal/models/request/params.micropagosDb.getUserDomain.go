package request

type GetUserRequest struct {
	UserName string
}

func NewGetUserRequest(userName string) *GetUserRequest {
	return &GetUserRequest{UserName: userName}
}
