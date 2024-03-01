package requests

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,validUsername"`
	Password string `json:"password" binding:"required,validPassword"`
	FullName string `json:"fullName" binding:"required,validFullname"`
	Email    string `json:"email" binding:"required,email"`
}

type GetUserRequest struct {
	Username string `uri:"username" binding:"required,validUsername"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,validUsername"`
	Password string `json:"password" binding:"required,validPassword"`
}
