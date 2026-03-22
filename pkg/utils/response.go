package utils

import "fmt"

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type StandardResponse struct {
	Success bool        `json:"success"`
	Data    any         `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
}

// 🔹 Asset Code Generator
func AssetCodeGen(year, seq int) string {
	return fmt.Sprintf("AST-%d-%04d", year, seq)
}