package store

import (
	"biatosh/entity"
	"biatosh/store/database"
	"crypto/sha256"
	"encoding/hex"
)

func convertUserToEntityUser(user *database.User) *entity.User {
	return &entity.User{
		ID:       int(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Name:     user.Name,
		Password: user.Password,
	}
}

func genHashPassword(password string) string {
	hashSha256Password := sha256.New()
	hashSha256Password.Write([]byte(password))
	return hex.EncodeToString(hashSha256Password.Sum(nil))
}
