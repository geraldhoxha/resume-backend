package service

import (
	"context"
	"strings"

	"github.com/geraldhoxha/resume-backend/config"
	"github.com/geraldhoxha/resume-backend/graph/model"
	"github.com/geraldhoxha/resume-backend/tools"

	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, input model.NewUser) (*model.User, error){
	db := config.GetDB()

	input.Password = tools.HashPassword(input.Password)

	user := model.User{
		ID: uuid.New().String(),
		Name: input.Name,
		Email: strings.ToLower(input.Email),
		Password: input.Password,
	}

	if err := db.Model(user).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UserGetByID(ctx context.Context, id string) (*model.User, error) {
	db := config.GetDB()

	var user model.User
	if err := db.Model(user).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	db := config.GetDB()

	var user model.User

	if err := db.Model(user).Where("email LIKE ?", email).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
