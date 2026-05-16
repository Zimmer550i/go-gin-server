package dto

type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required,min=1"`
}