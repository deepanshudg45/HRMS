package utils

import (
	"fmt"
	"strings"
)

func NormalizeText(value string) string {
	return strings.TrimSpace(value)
}

// AssetCodeGen generates an asset code from year and sequence number
// Format: YYYY + 6-digit sequence (e.g., 2026000001)
func AssetCodeGen(year, seq int) string {
	return fmt.Sprintf("%d%06d", year, seq)
}
