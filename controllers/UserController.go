package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

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
	userRepository.Redis = &infrastructures.Redis{}

	userService := new(services.UserService)
	userService.UserRepository = userRepository
	userService.Redis = &infrastructures.Redis{}

	userController := new(UserController)
	userController.UserService = userService
	userController.Redis = &infrastructures.Redis{}

	return userController
}

// UserController is
type UserController struct {
	UserService services.IUserService
	Redis       infrastructures.IRedis
}

// Login is
func (c *UserController) Login(res http.ResponseWriter, req *http.Request) {
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

		helpers.Response(res, http.StatusInternalServerError, &helpers.JSONResponse{
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

//UploadProfile get user profile
func (c *UserController) UploadProfile(res http.ResponseWriter, req *http.Request) {

	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDUser := fmt.Sprintf("%v", UserInfo["id"])

	intIDUser, err := strconv.Atoi(IDUser)
	if err != nil {
		resp := map[string]interface{}{
			"code":    5500,
			"message": "Error convert ID",
			"event":   "failed-store-file",
			"error":   err.Error(),
			"data":    nil,
		}

		logrus.WithFields(resp).Info(fmt.Sprintf("failed-store-file"))
		helpers.DirectResponse(res, http.StatusInternalServerError, resp)
		return
	}

	req.ParseMultipartForm(1000 << 20)

	fileTypeReq := req.Form.Get("type")
	fmt.Println("fileTyfileTypeReqpe")
	fmt.Println(fileTypeReq)

	if fileTypeReq != "sign_file" && fileTypeReq != "avatar" && fileTypeReq != "identity_file" && fileTypeReq != "sr_file" {
		resp := map[string]interface{}{
			"code":    4400,
			"message": "Error File Type Request",
			"event":   "validation-upload-profile",
			"error":   errors.New("failed validation"),
			"data":    nil,
		}

		logrus.WithFields(resp).Info(fmt.Sprintf("validation-upload-profile"))
		helpers.DirectResponse(res, http.StatusBadRequest, resp)
		return
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		resp := map[string]interface{}{
			"code":    5500,
			"message": "Error Get File",
			"event":   "failed-store-file",
			"error":   err.Error(),
			"data":    nil,
		}

		logrus.WithFields(resp).Info(fmt.Sprintf("failed-store-file"))
		helpers.DirectResponse(res, http.StatusInternalServerError, resp)
		return
	}

	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {

		resp := map[string]interface{}{
			"code":    5500,
			"message": "Error Validation File",
			"event":   "failed-store-file",
			"error":   err.Error(),
			"data":    nil,
		}

		logrus.WithFields(resp).Info(fmt.Sprintf("failed-store-file"))
		helpers.DirectResponse(res, http.StatusInternalServerError, resp)
		return

	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		resp := map[string]interface{}{
			"code":    4400,
			"message": "Only Image File (png,jpg,jpeg)",
			"event":   "validation-file",
			"error":   "",
			"data":    filetype,
		}

		logrus.WithFields(resp).Info(fmt.Sprintf("failed-store-file"))
		helpers.DirectResponse(res, http.StatusInternalServerError, resp)
		return

	}
	fmt.Println("filetype")
	fmt.Println(filetype)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		resp := map[string]interface{}{
			"code":    4400,
			"message": "Error seek file",
			"event":   "validation-file",
			"error":   err.Error(),
			"data":    nil,
		}

		logrus.WithFields(resp).Info(fmt.Sprintf("failed-store-file"))
		helpers.DirectResponse(res, http.StatusInternalServerError, resp)
		return
	}

	data, err := c.UserService.UpdateFile(req.Context(), file, handler.Filename, intIDUser, fileTypeReq)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"event": "update-profile-file",
			"nip":   fileTypeReq,
			"error": err,
		}).Info(fmt.Sprintf("failed-update-profile-file-%s", IDUser))

		helpers.Response(res, http.StatusInternalServerError, &helpers.JSONResponse{
			Code:    5500,
			Message: "Internal server error",
			Error:   err.Error(),
		})

		return
	}

	if data.Code != 2200 {
		logrus.WithFields(logrus.Fields{
			"code":  4400,
			"event": "update-profile-file",
			"nip":   fileTypeReq,
			"error": errors.New("not 2200"),
		}).Info(fmt.Sprintf("failed-update-profile-file-%s", IDUser))

	} else {
		logrus.WithFields(logrus.Fields{
			"code":   2200,
			"event":  "update-profile-file",
			"IDUser": fileTypeReq,
			"error":  nil,
		}).Info(fmt.Sprintf("success-update-profile-file-%s", IDUser))
	}

	helpers.DirectResponse(res, http.StatusOK, data)

	return
}

//GetAll get all users
func (c *UserController) GetAll(res http.ResponseWriter, req *http.Request) {

	users, err := c.UserService.GetAll(req.Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"event": "get-all-user",
			"error": err,
		}).Info(fmt.Sprintf("failed-get-all-user"))

		helpers.Response(res, http.StatusInternalServerError, &helpers.JSONResponse{
			Code:    5500,
			Message: "Internal server error",
			Error:   err.Error(),
		})
	}

	helpers.DirectResponse(res, http.StatusOK, users)
	return
}
