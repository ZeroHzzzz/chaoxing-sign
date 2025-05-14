package utils

import "math/rand"

func GenerateInviteCode() string {
	// 生成6位随机数字字母组合的邀请码
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
