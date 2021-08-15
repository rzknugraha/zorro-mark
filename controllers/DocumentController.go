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

// Initiate Variable
var validate *validator.Validate

// InitDocumentController is
func InitDocumentController() *DocumentController {

	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	documentService := services.InitDocumentService()

	return &DocumentController{
		DocumentService: documentService,
	}
}

// DocumentController is
type DocumentController struct {
	DocumentService services.IDocumentService
}

// GetDocuments is
func (c *DocumentController) GetDocuments(res http.ResponseWriter, req *http.Request) {
	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDUser := fmt.Sprintf("%v", UserInfo["id"])

	intIDUser, err := strconv.Atoi(IDUser)
	if err != nil {
		return
	}

	params := req.URL.Query()

	page, _ := strconv.Atoi(params.Get("page"))
	limit, _ := strconv.Atoi(params.Get("limit"))
	starred, _ := strconv.Atoi(params.Get("starred"))
	signing, _ := strconv.Atoi(params.Get("signing"))
	signed, _ := strconv.Atoi(params.Get("signed"))
	shared, _ := strconv.Atoi(params.Get("shared"))
	labels, _ := strconv.Atoi(params.Get("labels"))

	filter := models.DocumentUserFilter{
		Starred:  starred,
		Signing:  signing,
		Signed:   signed,
		Labels:   labels,
		Shared:   shared,
		UserID:   intIDUser,
		FileName: params.Get("file_name"),
		Sort:     params.Get("sort"),
	}

	pageReq := helpers.PageReq{
		Limit:  limit,
		Page:   page,
		Offset: (page - 1) * limit,
	}

	data, err := c.DocumentService.GetDocumentUser(req.Context(), filter, pageReq)
	if err != nil {
		return
	}

	helpers.DirectResponse(res, http.StatusOK, data)
	return

}

//UpdateDocument update document attribute
func (c *DocumentController) UpdateDocument(res http.ResponseWriter, req *http.Request) {

	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDUser := fmt.Sprintf("%v", UserInfo["id"])
	Name := fmt.Sprintf("%v", UserInfo["name"])
	NIP := fmt.Sprintf("%v", UserInfo["nip"])

	intIDUser, err := strconv.Atoi(IDUser)
	if err != nil {
		return
	}

	userData := models.Shortuser{
		Name: Name,
		Nip:  NIP,
		ID:   intIDUser,
	}

	var dataReq models.UpdateDocReq
	//Read request data
	body, _ := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(body, &dataReq)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataReq,
		}).Error("[CONTROLLER UpdateDocument] error parsing params")

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
		}).Error("[CONTROLLER UpdateDocument] Payload validation error")

		helpers.DirectResponse(res, http.StatusOK, &helpers.JSONResponse{
			Code:    4422,
			Message: "payload validation error",
			Error:   err.Error(),
			Data:    errFields,
		})
		return
	}

	dataReq.UserID = intIDUser

	result, err := c.DocumentService.UpdateDocumentAttributte(req.Context(), dataReq, userData)
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

// GetSingleDocument is
func (c *DocumentController) GetSingleDocument(res http.ResponseWriter, req *http.Request) {
	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDUser := fmt.Sprintf("%v", UserInfo["id"])
	Name := fmt.Sprintf("%v", UserInfo["name"])
	NIP := fmt.Sprintf("%v", UserInfo["nip"])

	intIDUser, err := strconv.Atoi(IDUser)
	if err != nil {
		return
	}

	userData := models.Shortuser{
		Name: Name,
		Nip:  NIP,
		ID:   intIDUser,
	}
	IDDoc := mux.Vars(req)["IDDoc"]

	intIDDoc, err := strconv.Atoi(IDDoc)
	if err != nil {
		return
	}

	data, err := c.DocumentService.GetSingleDocByUser(req.Context(), intIDDoc, userData)
	if err != nil {
		return
	}

	helpers.DirectResponse(res, http.StatusOK, data)
	return

}

// GetDocActivity is
func (c *DocumentController) GetDocActivity(res http.ResponseWriter, req *http.Request) {
	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDUser := fmt.Sprintf("%v", UserInfo["id"])
	Name := fmt.Sprintf("%v", UserInfo["name"])
	NIP := fmt.Sprintf("%v", UserInfo["nip"])

	intIDUser, err := strconv.Atoi(IDUser)
	if err != nil {
		return
	}

	userData := models.Shortuser{
		Name: Name,
		Nip:  NIP,
		ID:   intIDUser,
	}
	IDDoc := mux.Vars(req)["IDDoc"]

	intIDDoc, err := strconv.Atoi(IDDoc)
	if err != nil {
		return
	}

	data, err := c.DocumentService.GetActivityDoc(req.Context(), intIDDoc, userData)
	if err != nil {
		return
	}

	helpers.DirectResponse(res, http.StatusOK, data)
	return

}

//SaveDraft saving draft labels : 1
func (c *DocumentController) SaveDraft(res http.ResponseWriter, req *http.Request) {

	UserInfo, _ := gorillaContext.Get(req, "UserInfo").(jwt.MapClaims)

	IDUser := fmt.Sprintf("%v", UserInfo["id"])
	Name := fmt.Sprintf("%v", UserInfo["name"])
	NIP := fmt.Sprintf("%v", UserInfo["nip"])

	intIDUser, err := strconv.Atoi(IDUser)
	if err != nil {
		return
	}

	userData := models.Shortuser{
		Name: Name,
		Nip:  NIP,
		ID:   intIDUser,
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

	result, err := c.DocumentService.SaveDraft(req.Context(), userData, dataReq)
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
