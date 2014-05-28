package api

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"

	"paintToWin/lobby/user"
	"paintToWin/storage"
	"paintToWin/web"
)

const LoginResult_InvalidUserNameOrPassword = "InvalidUsernameOrPassword"

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (loginInput LoginInput) Validate() []web.InputError {
	return []web.InputError{}
}

type CreateUserInput struct {
	Email    string `json:"email" binding:"required"`
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (createUser CreateUserInput) Validate() []web.InputError {
	var errors []web.InputError
	if len(createUser.UserName) < 4 {
		errors = append(errors, web.NewInputError("userName", "User name must be at least 4 characters long"))
	} else if len(createUser.UserName) > 20 {
		errors = append(errors, web.NewInputError("userName", "User name cannot be longer than 20 characters"))
	}

	if matched, _ := regexp.MatchString(".*@.*\\..*", createUser.Email); !matched {
		errors = append(errors, web.NewInputError("email", "Not a valid email address"))
	}

	if len(createUser.Password) < 8 {
		errors = append(errors, web.NewInputError("password", "Password must be at least 8 characters long"))
	}
	return errors
}

func CreateUserHandler(store *storage.Storage) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		var input CreateUserInput
		if inputErrs, err := web.DeserializeAndValidateInput(req, &input); err != nil {
			return nil, web.NewApiError(http.StatusBadRequest, err.Error())
		} else if inputErrs != nil {
			return nil, web.NewApiError(http.StatusBadRequest, inputErrs)
		}

		err := user.CreateAccount(store, input.UserName, input.Email, input.Password)
		return nil, web.NewApiError(http.StatusInternalServerError, err.Error())
	}
}

func LoginHandler(store *storage.Storage) web.RequestHandler {
	return func(req *http.Request) (interface{}, web.ApiError) {
		var input LoginInput
		if inputErrs, err := web.DeserializeAndValidateInput(req, &input); err != nil {
			return nil, web.NewApiError(http.StatusBadRequest, err)
		} else if inputErrs != nil && len(inputErrs) > 0 {
			return nil, web.NewApiError(http.StatusBadRequest, inputErrs)
		}
		if userSession, err := user.Login(store, input.Email, input.Password); err != nil {
			return nil, web.NewApiError(http.StatusUnauthorized, LoginResult_InvalidUserNameOrPassword)
		} else {
			return NewSession(userSession.Id, userSession.Player.UserName), nil
		}
	}
}

func RegisterUserApi(router *mux.Router, store *storage.Storage) {
	router.HandleFunc("/users/create", web.DefaultHandler(CreateUserHandler(store))).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/login", web.DefaultHandler(LoginHandler(store))).Methods("POST", "OPTIONS")
}
