package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	fmt.Println(handler.Size)
	fmt.Println(handler.Filename)

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

	fmt.Println("buff")
	fmt.Println(buff)
	filetype := http.DetectContentType(buff)
	if filetype != "application/pdf" {
		resp := map[string]interface{}{
			"code":    4400,
			"message": "Only PDF File",
			"event":   "validation-file",
			"error":   "",
			"data":    nil,
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
	fmt.Println("file")
	fmt.Println(file)
	data, err := c.UploadService.StoreFile(req.Context(), file, handler.Filename)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":      5500,
			"event":     "failed-store-file",
			"file-info": handler,
			"error":     err,
		}).Error(fmt.Sprintf("failed-store-file"))

		helpers.DirectResponse(res, http.StatusInternalServerError, &helpers.JSONResponse{
			Code:    5500,
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	fmt.Println("data")
	fmt.Println(data)
	if data.Code != 2200 {
		logrus.WithFields(logrus.Fields{
			"code":      4400,
			"event":     "failed-store-file",
			"file-info": handler,
			"error":     "error store file",
		}).Error(fmt.Sprintf("failed-store-file"))

	} else {
		logrus.WithFields(logrus.Fields{
			"code":      2200,
			"event":     "success-store-file",
			"file-info": handler,
			"error":     nil,
		}).Info(fmt.Sprintf("failed-store-file"))
	}

	helpers.DirectResponse(res, http.StatusOK, data)
	return

}

// GetFile is
func (c *UploadController) GetFile(res http.ResponseWriter, req *http.Request) {

	fileName := mux.Vars(req)["filename"]

	path := viper.GetString("storage.path")
	// Open file
	f, err := os.Open("./" + path + "/" + fileName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":      5500,
			"event":     "failed-get-file",
			"file-info": fileName,
			"error":     err,
		}).Error(fmt.Sprintf("failed-get-file"))

		helpers.DirectResponse(res, http.StatusInternalServerError, &helpers.JSONResponse{
			Code:    5500,
			Message: "Internal server error",
			Error:   err.Error(),
			Data:    nil,
		})
		return
	}
	defer f.Close()

	//Set header
	res.Header().Set("Content-type", "application/pdf")

	//Stream to response
	if _, err := io.Copy(res, f); err != nil {
		logrus.WithFields(logrus.Fields{
			"code":      5500,
			"event":     "failed-stream-file",
			"file-info": fileName,
			"error":     err,
		}).Error(fmt.Sprintf("failed-stream-file"))

		helpers.DirectResponse(res, http.StatusInternalServerError, &helpers.JSONResponse{
			Code:    5500,
			Message: "Internal server error",
			Error:   err.Error(),
			Data:    nil,
		})
		return
	}

	return

}
