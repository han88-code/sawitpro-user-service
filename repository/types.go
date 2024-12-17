// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	FName string
}

type User struct {
	Id           int    `json:"id"`
	Fullname     string `json:"fullname"`
	Phonenumber  string `json:"phonenumber"`
	Password     string `json:"password"`
	Successlogin int    `json:"successlogin"`
}
