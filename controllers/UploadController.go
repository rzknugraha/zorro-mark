package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	gorillaContext "github.com/gorilla/context"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/services"
	"github.com/sirupsen/logrus"
)

// InitUploadController is
func InitUploadController() *UploadController {

	uploadService := services.InitUploadService()

	return &UploadController{
		UploadService: uploadService,
	}
}

// UploadController is
type UploadController struct {
	UploadService services.IUploadService
}

// Upload is
func (c *UploadController) Upload(res http.ResponseWriter, req *http.Request) {
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

	data, err := c.UploadService.StoreFile(req.Context(), file, handler.Filename, intIDUser)
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

	var path models.FileReq
	//Read request data
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &path)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  path,
		}).Error("[CONTROLLER GetFile] error parsing params")

		helpers.DirectResponse(res, http.StatusBadRequest, "Failed read input data")
		return
	}
	if path.Path == "" {
		logrus.WithFields(logrus.Fields{
			"code":      4400,
			"event":     "failed-get-file",
			"file-info": "",
			"error":     err,
		}).Error(fmt.Sprintf("[CONTROLLER GetFile] failed-get-file"))

		helpers.DirectResponse(res, http.StatusInternalServerError, &helpers.JSONResponse{
			Code:    4400,
			Message: "Error or empty path name",
			Error:   "true",
			Data:    nil,
		})
		return

	}
	// Open file
	f, err := os.Open("." + path.Path)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":      5500,
			"event":     "failed-get-file",
			"file-info": path,
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
			"file-info": path,
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
