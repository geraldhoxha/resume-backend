package service

import (
	"context"

	"github.com/geraldhoxha/resume-backend/graph/model"
	"github.com/geraldhoxha/resume-backend/tools"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

func UserRegister(ctx context.Context, input model.NewUser) (*model.AuthResponse, error){
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

	token, err := JwtGenerate(ctx, createdUser.ID, createdUser.Name, createdUser.Email)
	if err != nil {
		return nil, err
	}

	response := &model.AuthResponse{
		Token: token,
		User: createdUser,
	}

	return response, nil
}

func UserLogin(ctx context.Context, email string, password string) (*model.AuthResponse, error){
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

	token, err := JwtGenerate(ctx, getUser.ID, getUser.Name, getUser.Email)
	if err != nil {
		return nil, err
	}
	
	response := &model.AuthResponse{
		Token: token,
		User: getUser,
	}

	return response ,nil
}
