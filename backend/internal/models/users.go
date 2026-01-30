package models

type CreateUserRequest struct {
	Handle        string `json:"handle" binding:"required,min=4,max=16"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
	PasswordCheck string `json:"password_check" binding:"required,eqfield=Password"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserToken struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type CreateDMRequest struct {
	RecipientID int `json:"recipient_id" binding:"required"`
}
