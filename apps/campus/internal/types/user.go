package types

type RegisterUser struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
