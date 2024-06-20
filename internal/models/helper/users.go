package helper

import (
	"context"
	"errors"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

var Users []response.UserInfo

func GetUserdata(ctx context.Context, username string) (response.UserResponse, error) {
	ctx = ins_log.SetPackageNameInContext(ctx, "handler")
	ins_log.Infof(ctx, "starting to get the userdata for the user %s", username)
	userData := response.UserResponse{}
	for _, user := range Users {
		if user.Username == username {
			userData.Username = user.Username
			userData.Password = user.UserPassword
			userData.UserDomain = user.DomainName
			userData.Result = 0
			userData.Message = "the domain name is: " + user.DomainName
			return userData, nil
		}
	}
	return userData, errors.New("user not found in the list of users")
}

func UpdatedUserLastLogin(ctx context.Context, username string) {

}
