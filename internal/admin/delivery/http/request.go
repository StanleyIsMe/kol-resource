package http

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	UserName string `json:"user_name" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	AdminName string `json:"admin_name"`
}
