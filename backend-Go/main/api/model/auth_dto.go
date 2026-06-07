package model

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UserOrEmail string `json:"userOrEmail"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
