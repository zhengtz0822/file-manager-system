package uuid

import (
	"github.com/google/uuid"
)

// Generate 生成 UUID
func Generate() string {
	return uuid.New().String()
}
