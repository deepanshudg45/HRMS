package model

type LoginRequest struct {
	EmployeeID string `json:"employeeId"`
	Role       string `json:"role"`
}

type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	EmployeeID string `json:"employeeId"`
	Role       string `json:"role"`
}
