package handler

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/util"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Function for check register input value.
func RegisterInputValidation(input generated.RegisterRequestBody) (string, bool) {
	errorMessage := ""
	isValid := true

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	// Check Phone Number length
	if len(input.Phonenumber) < 10 || len(input.Phonenumber) > 13 {
		errorMessage += "Phone Number length must between 10 and 13; "
	}

	// Check Phone Number contain +62
	isPhonePrefixValid := strings.HasPrefix(input.Phonenumber, "+62")
	if !isPhonePrefixValid {
		errorMessage += "Phone Number must start with +62; "
	}

	// Check Full Name length
	if len(input.Fullname) < 3 || len(input.Fullname) > 60 {
		errorMessage += "Full Name length must between 3 and 60; "
	}

	// Check Password length
	if len(input.Password) < 6 || len(input.Password) > 64 {
		errorMessage += "Full Name length must between 6 and 64; "
	}

	// Check Password characters
	for _, char := range input.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		errorMessage += "Password must contain at least 1 capital characters and 1 number and 1 special (non alpha-numeric) characters;"
	}

	if errorMessage != "" {
		isValid = false
	}

	return errorMessage, isValid
}

// Function for hashing password.
func HashPassword(password string) (string, error) {
	bytesPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytesPassword), err
}

// Function for Register User.
func (s *Server) Register(ctx echo.Context) error {

	u := new(generated.RegisterRequestBody)
	if err := ctx.Bind(u); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	userInput := generated.RegisterRequestBody{
		Fullname:    u.Fullname,
		Phonenumber: u.Phonenumber,
		Password:    u.Password,
	}

	// Check register input validity
	errorMessage, isValid := RegisterInputValidation(userInput)
	if !isValid {
		var resp generated.ErrorResponse
		resp.Message = fmt.Sprintf("Error: %s", errorMessage)

		return ctx.JSON(http.StatusBadRequest, resp)
	}

	// Hashing password
	hashPassword, errPass := HashPassword(userInput.Password)
	if errPass != nil {
		return ctx.JSON(http.StatusBadRequest, errPass)
	}

	userInput.Password = hashPassword

	// Insert new user to database
	newUser, errDb := s.Repository.InsertUser(ctx.Request().Context(), userInput)
	if errDb != nil {
		return ctx.JSON(http.StatusBadRequest, errDb)
	}

	var resp generated.RegisterSuccessResponse
	resp.Id = newUser.Id

	return ctx.JSON(http.StatusOK, resp)
}

// Function for Login User.
func (s *Server) Login(ctx echo.Context) error {
	u := new(generated.LoginRequestBody)
	if err := ctx.Bind(u); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	userInput := generated.LoginRequestBody{
		Password:    u.Password,
		Phonenumber: u.Phonenumber,
	}

	// Hashing password
	hashPassword, errPass := HashPassword(userInput.Password)
	if errPass != nil {
		return ctx.JSON(http.StatusBadRequest, errPass)
	}

	userInput.Password = hashPassword

	// get user data base on username & password
	user, errDb := s.Repository.GetUser(ctx.Request().Context(), userInput)
	if errDb != nil {
		return ctx.JSON(http.StatusBadRequest, errDb)
	}

	// increasing number of successful login
	user.Successlogin++

	// update user number of successful login
	_, errDbUpdate := s.Repository.UpdateUserSuccesslogin(ctx.Request().Context(), user)
	if errDbUpdate != nil {
		return ctx.JSON(http.StatusBadRequest, errDbUpdate)
	}

	// generate new token
	token, errToken := util.GenerateRSAToken(uint(user.Id))
	if errToken != nil {
		return ctx.JSON(http.StatusBadRequest, errToken)
	}

	var resp generated.LoginSuccessResponse
	resp.Id = user.Id
	resp.Token = token

	return ctx.JSON(http.StatusOK, resp)
}

// Function for Get User Data.
func (s *Server) Profile(ctx echo.Context, params generated.ProfileParams) error {
	userInput := generated.ProfileParams{
		Id: params.Id,
	}

	// get user data base on user id
	user, errDb := s.Repository.GetUserById(ctx.Request().Context(), userInput)
	if errDb != nil {
		return ctx.JSON(http.StatusForbidden, errDb)
	}

	return ctx.JSON(http.StatusOK, user)
}

// Function for Update User Data.
func (s *Server) Profileupdate(ctx echo.Context) error {
	u := new(generated.ProfileUpdateRequestBody)
	if err := ctx.Bind(u); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	userIdInput := generated.ProfileParams{
		Id: u.Id,
	}

	// get user data base on user id
	user, errDb := s.Repository.GetUserById(ctx.Request().Context(), userIdInput)
	if errDb != nil {
		return ctx.JSON(http.StatusForbidden, errDb)
	}

	userInput := generated.ProfileUpdateRequestBody{
		Id:          u.Id,
		Fullname:    u.Fullname,
		Phonenumber: u.Phonenumber,
	}

	if *userInput.Phonenumber != "" {
		// get user data base on phone number, check whether the updated phone number has been using or not.
		_, errDbOther := s.Repository.GetUserByPhonenumber(ctx.Request().Context(), userInput)
		if errDbOther != nil {
			return ctx.JSON(http.StatusConflict, errDbOther)
		}
	} else if *userInput.Phonenumber == "" {
		userInput.Phonenumber = &user.Phonenumber
	}

	if *userInput.Fullname == "" {
		userInput.Fullname = &user.Fullname
	}

	// update user full name & user phone number
	_, errDbUpdate := s.Repository.UpdateUser(ctx.Request().Context(), userInput)
	if errDbUpdate != nil {
		return ctx.JSON(http.StatusForbidden, errDbUpdate)
	}

	return ctx.JSON(http.StatusOK, userInput)
}
