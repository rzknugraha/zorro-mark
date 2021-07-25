package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"time"

	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
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

	// _ = writer.WriteField("nik", dataSign.NIK)
	_ = writer.WriteField("nik", "0803202100007062")
	_ = writer.WriteField("passphrase", dataSign.Passphrase)
	_ = writer.WriteField("tampilan", dataSign.Tampilan)
	// _ = writer.WriteField("page", dataSign.Page)
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

	fmt.Println("part1")
	fmt.Println(part1)
	// application/pdf

	req, err := http.NewRequest("POST", "http://192.168.1.31/api/sign/pdf", payload)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error post data")
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.SetBasicAuth("admin", "qwerty")

	fmt.Println("req")
	fmt.Println(req)

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
	body, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println("response Body:", string(body))

	result.StatusCode = rsp.StatusCode
	result.Message = string(body)

	if rsp.StatusCode != http.StatusOK {

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

	}
	return

}
