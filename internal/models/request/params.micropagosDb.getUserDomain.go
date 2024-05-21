package request

type GetUserDomainRequest struct {
	UserName string
}

func NewGetUserDomainRequest(userName string) *GetUserDomainRequest {
	return &GetUserDomainRequest{UserName: userName}
}
