package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/rzknugraha/zorro-mark/services"
)

// InitUserController is
func InitUserController() *UserController {
	userRepository := new(repositories.UserRepository)
	userRepository.DB = &infrastructures.SQLConnection{}

	userService := new(services.UserService)
	userService.UserRepository = userRepository

	userController := new(UserController)
	userController.UserService = userService

	return userController
}

// UserController is
type UserController struct {
	UserService services.IUserService
}

// Login is
func (c *UserController) Login(res http.ResponseWriter, req *http.Request) {
	var l models.Login
	//Read request data
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &l)

	if err != nil {
		helpers.Response(res, http.StatusBadRequest, "Failed read input data")
		return
	}

	result, err := c.UserService.Login(l)
	if err == nil {
		resp := map[string]interface{}{
			"code":    2200,
			"message": "",
			"error":   false,
			"data":    result,
		}
		helpers.Response(res, http.StatusOK, resp)
	} else {
		resp := map[string]interface{}{
			"code":    4400,
			"message": "Wrong Username or Passowrd",
			"error":   true,
			"data":    result,
		}
		helpers.Response(res, http.StatusBadRequest, resp)
	}

	return
}
