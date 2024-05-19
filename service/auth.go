package service

import (
	"context"

	"github.com/geraldhoxha/resume-backend/graph/model"
	"github.com/geraldhoxha/resume-backend/tools"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func UserRegister(ctx context.Context, input model.NewUser) (interface{}, error){
	_, err := UserGetByEmail(ctx, input.Email)
	if err != nil {

		if err != gorm.ErrRecordNotFound{
			return nil, err
		}
	}

	createdUser, err := CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	token, err := JwtGenerate(ctx, createdUser.ID, createdUser.Name)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

func UserLogin(ctx context.Context, email string, password string) (interface{}, error){
	getUser, err := UserGetByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &gqlerror.Error {
				Message: "Email not found",
			}
		}
		return nil, err
	}

	if err := tools.ComparePassword(getUser.Password, password); err != nil {
		return nil, err
	}

	token, err := JwtGenerate(ctx, getUser.ID, getUser.Name)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	},nil
}
