package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
)

func (r *Repository) InsertUser(ctx context.Context, input generated.RegisterRequestBody) (output generated.RegisterSuccessResponse, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO user (fullname,phonenumber,password) VALUES (?,?,?)", input.Fullname, input.Phonenumber, input.Password).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUser(ctx context.Context, input generated.LoginRequestBody) (output User, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT * FROM user WHERE phonenumber = ? AND password = ?", input.Phonenumber, input.Password).Scan(&output)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserById(ctx context.Context, input generated.ProfileParams) (output generated.ProfileSuccessResponse, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT fullname, phonenumber FROM user WHERE id = ?", input.Id).Scan(&output)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhonenumber(ctx context.Context, input generated.ProfileUpdateRequestBody) (output generated.ProfileSuccessResponse, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT fullname FROM user WHERE phonenumber = ?", input.Phonenumber).Scan(&output.Fullname)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateUserSuccesslogin(ctx context.Context, input User) (output generated.ProfileSuccessResponse, err error) {
	err = r.Db.QueryRowContext(ctx, "UPDATE user SET successlogin = ? WHERE id = ?", input.Successlogin, input.Id).Scan(&output)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateUser(ctx context.Context, input generated.ProfileUpdateRequestBody) (output generated.ProfileSuccessResponse, err error) {
	err = r.Db.QueryRowContext(ctx, "UPDATE user SET fullname = ?, phonenumber = ? WHERE id = ?", input.Fullname, input.Phonenumber, input.Id).Scan(&output)
	if err != nil {
		return
	}
	return
}
