package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IEsignRepository is
type IEsignRepository interface {
	PostEsign(ctx context.Context, dataSign models.EsignReq) (result models.EsignResp, err error)
}

// EsignRepository is
type EsignRepository struct {
	DB infrastructures.ISQLConnection
}

// PostEsign post to bsre
func (r *EsignRepository) PostEsign(ctx context.Context, dataSign models.EsignReq) (result models.EsignResp, err error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, errFile1 := os.Open("." + dataSign.FilePath)

	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": errFile1,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error get file stats")
		return
	}

	path := viper.GetString("storage.path")
	fileResp := models.UploadResp{
		FileName: fmt.Sprintf("signed-%s", fi.Name()),
	}
	fullPath := fmt.Sprintf("/%s/signed/%s", path, fileResp.FileName)

	partHeader := textproto.MIMEHeader{}
	disp := fmt.Sprintf("form-data; name=file; filename=%s", fi.Name())
	partHeader.Add("Content-Disposition", disp)
	partHeader.Add("Content-Type", "application/pdf")
	part1, errFile1 := writer.CreatePart(partHeader)

	// part1, errFile1 := writer.CreateFormFile("file", fi.Name())
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": errFile1,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error get file not founf")
		return
	}

	//using QR or Image
	if dataSign.Tampilan == "visible" {
		strPage := strconv.Itoa(dataSign.Page)
		strXAxis := strconv.Itoa(dataSign.XAxis)
		strYAxis := strconv.Itoa(dataSign.YAxis)
		strWidth := strconv.Itoa(dataSign.Width)
		strHeight := strconv.Itoa(dataSign.Height)

		_ = writer.WriteField("page", strPage)

		_ = writer.WriteField("xAxis", strXAxis)
		_ = writer.WriteField("yAxis", strYAxis)
		_ = writer.WriteField("width", strWidth)
		_ = writer.WriteField("height", strHeight)

		if dataSign.Image == true {
			_ = writer.WriteField("image", "true")

			fileSign, errFile2 := os.Open("." + dataSign.ImagePath)

			defer fileSign.Close()

			fi1, errFile2 := fileSign.Stat()
			if errFile2 != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": errFile2,
					"data":  dataSign,
				}).Error("[REPO PostEsign] error get file stats sign")
				return
			}

			fileExtension := filepath.Ext(fi1.Name())

			mimeCt := "image/jpeg"

			if fileExtension == ".png" {
				mimeCt = "image/png"
			}

			partHeader1 := textproto.MIMEHeader{}
			disp1 := fmt.Sprintf("form-data; name=imageTTD; filename=%s", fi1.Name())
			partHeader1.Add("Content-Disposition", disp1)
			partHeader1.Add("Content-Type", mimeCt)
			part2, errFile2 := writer.CreatePart(partHeader1)

			_, errFile2 = io.Copy(part2, fileSign)
			if errFile2 != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": errFile2,
					"data":  dataSign,
				}).Error("[REPO PostEsign] error get file sign not found")
				return
			}

		} else {
			_ = writer.WriteField("image", "false")
			_ = writer.WriteField("linkQR", viper.GetString("static_file")+fullPath)
		}
	}

	if viper.GetString("esign.env") == "prod" {
		_ = writer.WriteField("nik", dataSign.NIK)
	} else {
		_ = writer.WriteField("nik", viper.GetString("esign.identity_no"))
	}

	_ = writer.WriteField("passphrase", dataSign.Passphrase)
	_ = writer.WriteField("tampilan", dataSign.Tampilan)
	//
	// _ = writer.WriteField("image", dataSign.Image)
	// _ = writer.WriteField("width", dataSign.Width)
	// _ = writer.WriteField("height", dataSign.Height)
	err = writer.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error close writter")
		return
	}

	// application/pdf

	req, err := http.NewRequest("POST", viper.GetString("esign.url"), payload)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error post data")
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.SetBasicAuth(viper.GetString("esign.username"), viper.GetString("esign.password"))

	rsp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error make client do")
		return
	}
	fmt.Println("rsp")
	fmt.Println(rsp)

	defer rsp.Body.Close()

	fmt.Println("response Status:", rsp.Status)
	fmt.Println("response Headers:", rsp.Header)

	if rsp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(rsp.Body)
		err = json.Unmarshal(body, &result)
		if err != nil {

			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  rsp.Body,
			}).Error("[REPO PostEsign] error unmarshall ")

			return
		}
		if rsp.StatusCode == http.StatusBadRequest {

			logrus.WithFields(logrus.Fields{
				"code":  4400,
				"error": err,
				"data":  rsp.Body,
			}).Error("[REPO PostEsign] error make client do")

		} else {

			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  rsp.Body,
			}).Error("[REPO PostEsign] error make client do")
		}

		return

	}

	dst, err := os.Create("." + fullPath)
	if err != nil {
		return

	}

	defer dst.Close()
	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, rsp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  rsp.Body,
		}).Error("[REPO PostEsign] error move file signed")
		return
	}

	result.StatusCode = 200
	result.PathFile = fullPath

	return

}
