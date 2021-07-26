package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	gorillaContext "github.com/gorilla/context"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/services"
	"github.com/sirupsen/logrus"
)

// InitEsignController is
func InitEsignController() *EsignController {

	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	esignService := services.InitEsignService()

	return &EsignController{
		EsignService: esignService,
	}
}

// EsignController is
type EsignController struct {
	EsignService services.IEsignService
}

//SignDoc update esign attribute
func (c *EsignController) SignDoc(res http.ResponseWriter, req *http.Request) {

	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDUser := fmt.Sprintf("%v", UserInfo["id"])
	Name := fmt.Sprintf("%v", UserInfo["name"])
	NIP := fmt.Sprintf("%v", UserInfo["nip"])
	SignFile := fmt.Sprintf("%v", UserInfo["sign_file"])
	IdentityNO := fmt.Sprintf("%v", UserInfo["identity_no"])

	intIDUser, err := strconv.Atoi(IDUser)
	if err != nil {
		return
	}

	userData := models.Shortuser{
		Name:       Name,
		Nip:        NIP,
		ID:         intIDUser,
		IdentityNO: IdentityNO,
		SignFile:   SignFile,
	}

	fmt.Println(userData)
	var dataReq models.EsignReq
	//Read request data
	body, _ := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(body, &dataReq)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataReq,
		}).Error("[CONTROLLER SignDoc] error parsing params")

		helpers.DirectResponse(res, http.StatusBadRequest, "Failed read input data")
		return
	}

	if err := validate.Struct(dataReq); err != nil {
		errField := map[string]string{}
		errFields := []map[string]string{}

		for _, e := range err.(validator.ValidationErrors) {
			errField[e.Field()] = fmt.Sprintf("%s failed on the %s tag", e.Field(), e.Tag())
		}
		errFields = append(errFields, errField)

		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataReq,
		}).Error("[CONTROLLER SignDoc] Payload validation error")

		helpers.DirectResponse(res, http.StatusOK, &helpers.JSONResponse{
			Code:    4422,
			Message: "payload validation error",
			Error:   err.Error(),
			Data:    errFields,
		})
		return
	}
	dataReq.NIK = userData.IdentityNO
	dataReq.ImagePath = userData.SignFile

	result, err := c.EsignService.PostSign(req.Context(), dataReq, userData)
	if err != nil {
		responseErr := &helpers.JSONResponse{
			Code:    5500,
			Message: "Error Internal",
			Data:    nil,
		}

		helpers.DirectResponse(res, http.StatusInternalServerError, responseErr)
		return
	}

	helpers.DirectResponse(res, http.StatusOK, result)

	return
}
