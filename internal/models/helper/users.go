package helper

import (
	"context"
	"errors"
	"time"

	"github.com/scch94/MICROPAGOSMESSAGEGATEWAY/internal/models/response"
	"github.com/scch94/ins_log"
)

var Users []response.UserInfo

var Userslastlogin []response.UpdateLastLogin

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

	ctx = ins_log.SetPackageNameInContext(ctx, "handler")
	ins_log.Tracef(ctx, "starting to updated the last login of the user %s", username)
	time.Sleep(1 * time.Minute)
	// Get the current time
	now := time.Now()

	// Format the time as a string
	formattedTime := now.Format("2006-01-02 15:04:05")

	for _, user := range Userslastlogin {
		if user.UserName == username {
			user.LoginTime = formattedTime
			ins_log.Tracef(ctx, "user last login updated to: %s", user.LoginTime)
			return
		}
	}
	Userslastlogin = append(Userslastlogin, response.UpdateLastLogin{UserName: username, LoginTime: formattedTime})
	ins_log.Tracef(ctx, "user %s inserted in the lastlogin data", username)
}
