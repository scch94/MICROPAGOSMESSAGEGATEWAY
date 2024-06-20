package response

type UpdateUsersLastLoginResponse struct {
	Result       uint8  `json:"result"`
	Message      string `json:"message"`
	RowsAffected int64  `json:"rowsaffected"`
}
