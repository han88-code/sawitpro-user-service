// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
)

type RepositoryInterface interface {
	InsertUser(ctx context.Context, input generated.RegisterRequestBody) (output generated.RegisterSuccessResponse, err error)
	GetUser(ctx context.Context, input generated.LoginRequestBody) (output User, err error)
	GetUserById(ctx context.Context, input generated.ProfileParams) (output generated.ProfileSuccessResponse, err error)
	GetUserByPhonenumber(ctx context.Context, input generated.ProfileUpdateRequestBody) (output generated.ProfileSuccessResponse, err error)
	UpdateUserSuccesslogin(ctx context.Context, input User) (output generated.ProfileSuccessResponse, err error)
	UpdateUser(ctx context.Context, input generated.ProfileUpdateRequestBody) (output generated.ProfileSuccessResponse, err error)
}
