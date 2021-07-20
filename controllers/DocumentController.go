package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	gorillaContext "github.com/gorilla/context"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/services"
)

// InitDocumentController is
func InitDocumentController() *DocumentController {

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
