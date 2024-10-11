package http

import (
	"fmt"

	"github.com/goccy/go-json"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	UserName string `json:"user_name" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (r *RegisterRequest) String() string {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("marshalling RegisterRequest error: %v", err)
	}

	return string(jsonStr)
}

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	AdminName string `json:"admin_name"`
}
