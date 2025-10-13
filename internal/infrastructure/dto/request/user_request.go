package request

type CreateUser struct {
	Email    string `json:"email" binding:"required,max=100"`
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required"`
}

type LogIn struct {
	EmailOrUsername string `json:"emailOrUsername" binding:"required,max=100"`
	Password        string `json:"password" binding:"required"`
}
