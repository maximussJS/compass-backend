package password

import (
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+=-"

func GenerateRandomPassword() string {
	var password []byte

	for i := 0; i < 16; i++ {
		randNum := rand.Intn(len(charset))
		password = append(password, charset[randNum])
	}
	return string(password)
}
