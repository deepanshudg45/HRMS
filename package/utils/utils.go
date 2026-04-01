package utils

import (
	"fmt"
)

func AssetCodeGen(year, seq int) string {
	return fmt.Sprintf("%d%06d", year, seq)
}
