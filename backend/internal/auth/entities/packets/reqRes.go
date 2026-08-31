package packets

type UserDataReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDataResponse struct {
	UserID   string
	Username string
	Email    string
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  UserDataResponse
	Token string
}
