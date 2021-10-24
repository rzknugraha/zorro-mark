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
	"github.com/gorilla/mux"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/services"
	"github.com/sirupsen/logrus"
)

// InitCommentController is
func InitCommentController() *CommentController {

	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	commentService := services.InitCommentService()

	return &CommentController{
		CommentService: commentService,
	}
}

// CommentController is
type CommentController struct {
	CommentService services.ICommentService
}

// GetComments is
func (c *CommentController) GetComments(res http.ResponseWriter, req *http.Request) {
	// UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDDoc := mux.Vars(req)["IDDoc"]

	intIDDoc, err := strconv.Atoi(IDDoc)
	if err != nil {
		return
	}

	data, err := c.CommentService.GetCommentByDocID(req.Context(), intIDDoc)
	if err != nil {
		helpers.DirectResponse(res, http.StatusInternalServerError, data)

		return
	}

	helpers.DirectResponse(res, http.StatusOK, data)
	return

}

// StoreComment is
func (c *CommentController) StoreComment(res http.ResponseWriter, req *http.Request) {
	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDUser := fmt.Sprintf("%v", UserInfo["id"])
	Name := fmt.Sprintf("%v", UserInfo["name"])

	intIDUser, err := strconv.Atoi(IDUser)
	if err != nil {
		return
	}

	var dataReq models.Comment
	//Read request data
	body, _ := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(body, &dataReq)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataReq,
		}).Error("[CONTROLLER StoreComment] error parsing params")

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
		}).Error("[CONTROLLER StoreComment] Payload validation error")

		helpers.DirectResponse(res, http.StatusOK, &helpers.JSONResponse{
			Code:    4422,
			Message: "payload validation error",
			Error:   err.Error(),
			Data:    errFields,
		})
		return
	}

	dataReq.NameUser = Name

	data, err := c.CommentService.StoreComment(req.Context(), intIDUser, dataReq)
	if err != nil {
		helpers.DirectResponse(res, http.StatusInternalServerError, data)

		return
	}

	helpers.DirectResponse(res, http.StatusOK, data)
	return

}
