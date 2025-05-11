package verification

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CodeLength = 6
	CodeTTL    = 5 * time.Minute
)

func GenerateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := r.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}

func VerifyCode(storedCode, inputCode string) bool {
	return storedCode == inputCode
}
