// internal/auth/dto.go
package auth

type RegisterRequest struct {
	WorkEmail    string `json:"work_email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	EmployeeCode string `json:"employee_code" binding:"required"`
}

type LoginRequest struct {
	WorkEmail string `json:"work_email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}
