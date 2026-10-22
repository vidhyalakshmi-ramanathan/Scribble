package packets

type UserDataReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDataResponse struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  UserDataResponse `json:"user"`
	Token string           `json:"token"`
}
