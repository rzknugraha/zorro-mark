package controllers

import (
	"fmt"
	"net/http"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/services"
	"github.com/sirupsen/logrus"
)

// InitUploadController is
func InitUploadController() *UploadController {

	uploadService := new(services.UploadService)

	uploadController := new(UploadController)
	uploadController.UploadService = uploadService

	return uploadController
}

// UploadController is
type UploadController struct {
	UploadService services.IUploadService
}

// Upload is
func (c *UploadController) Upload(res http.ResponseWriter, req *http.Request) {

	req.ParseMultipartForm(10 << 20)

	file, handler, err := req.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	fmt.Println(handler.Size)

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
		helpers.Response(res, http.StatusInternalServerError, resp)
		return

	}

	filetype := http.DetectContentType(buff)
	if filetype != "application/pdf" {
		resp := map[string]interface{}{
			"code":    4400,
			"message": "Only PDF File",
			"event":   "validation-file",
			"error":   err.Error(),
			"data":    nil,
		}

		logrus.WithFields(resp).Info(fmt.Sprintf("failed-store-file"))
		helpers.Response(res, http.StatusInternalServerError, resp)
		return

	}

	data, err := c.UploadService.StoreFile(req.Context(), file,handler.Filename)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":      5500,
			"event":     "failed-store-file",
			"file-info": handler,
			"error":     err,
		}).Error(fmt.Sprintf("failed-store-file"))

		helpers.Response(res, http.StatusOK, &helpers.JSONResponse{
			Code:    5500,
			Message: "Internal server error",
			Error:   err.Error(),
		})
	}

	if data.Code != 2200 {
		logrus.WithFields(logrus.Fields{
			"code":      4400,
			"event":     "failed-store-file",
			"file-info": handler,
			"error":     err,
		}).Error(fmt.Sprintf("failed-store-file"))

	} else {
		logrus.WithFields(logrus.Fields{
			"code":      2200,
			"event":     "success-store-file",
			"file-info": handler,
			"error":     err,
		}).Info(fmt.Sprintf("failed-store-file"))
	}

	helpers.Response(res, http.StatusOK, data)
	return

}
