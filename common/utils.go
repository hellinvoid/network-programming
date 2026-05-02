package common

import (
	"unicode"

	"github.com/google/uuid"
)

func IsAlnum(s string) bool {
    for _, r := range s {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
            return false
        }
    }
    return true
}
func GenerateUUID() string {
    return uuid.New().String()
}