package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	gorillaContext "github.com/gorilla/context"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/rzknugraha/zorro-mark/services"
	"github.com/sirupsen/logrus"
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

	//Allow CORS here By * or specific origin
	res.Header().Set("Access-Control-Allow-Origin", "*")

	var l models.Login
	//Read request data
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &l)

	if err != nil {
		helpers.DirectResponse(res, http.StatusBadRequest, "Failed read input data")
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
		helpers.DirectResponse(res, http.StatusOK, resp)
	} else {
		resp := map[string]interface{}{
			"code":    4400,
			"message": "Wrong Username or Passowrd",
			"error":   true,
			"data":    result,
		}
		helpers.DirectResponse(res, http.StatusBadRequest, resp)
	}

	return
}

//Profile get user profile
func (c *UserController) Profile(res http.ResponseWriter, req *http.Request) {

	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	nip := fmt.Sprintf("%s", UserInfo["nip"])

	if len(nip) == 0 {
		resp := map[string]interface{}{
			"code":    4400,
			"message": "Username Not Found",
			"error":   true,
			"data":    nil,
		}
		logrus.WithFields(resp).Error("[Profile Controller] get context nil")
		helpers.DirectResponse(res, http.StatusBadRequest, resp)
		return
	}

	fmt.Println("nip")
	fmt.Println(nip)

	data, err := c.UserService.Profile(req.Context(), nip)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"event": "get-profile",
			"nip":   nip,
			"error": err,
		}).Info(fmt.Sprintf("failed-get-profile-%s", nip))

		helpers.Response(res, http.StatusOK, &helpers.JSONResponse{
			Code:    5500,
			Message: "Internal server error",
			Error:   err.Error(),
		})
	}

	if data.Code != 2200 {
		logrus.WithFields(logrus.Fields{
			"code":  4400,
			"event": "get-profile",
			"nip":   nip,
			"error": errors.New("not 2200"),
		}).Info(fmt.Sprintf("failed-get-profile-%s", nip))

	} else {
		logrus.WithFields(logrus.Fields{
			"code":  2200,
			"event": "get-profile",
			"nip":   nip,
			"error": nil,
		}).Info(fmt.Sprintf("success-get-profile-%s", nip))
	}

	helpers.DirectResponse(res, http.StatusOK, data)
	return
}
