package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
)

// IEsignRepository is
type IEsignRepository interface {
	PostEsign(ctx context.Context, dataSign models.EsignReq) (err error)
}

// EsignRepository is
type EsignRepository struct {
	DB infrastructures.ISQLConnection
}

// PostEsign post to bsre
func (r *EsignRepository) PostEsign(ctx context.Context, dataSign models.EsignReq) (err error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, errFile1 := os.Open("." + dataSign.FilePath)
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("file", filepath.Base(dataSign.FilePath))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": errFile1,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error get file not founf")
		return
	}

	_ = writer.WriteField("nik", dataSign.NIK)
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
	req, err := http.NewRequest("POST", "http://192.168.1.31/api/sign/pdf", payload)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error post data")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.SetBasicAuth("admin", "qwerty")

	rsp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataSign,
		}).Error("[REPO PostEsign] error make client do")
		return err
	}
	fmt.Println("rsp")
	fmt.Println(rsp)

	fmt.Println("req")
	fmt.Println(req)

	defer rsp.Body.Close()

	fmt.Println("response Status:", rsp.Status)
	fmt.Println("response Headers:", rsp.Header)
	body, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println("response Body:", string(body))

	if rsp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  rsp.Body,
		}).Error("[REPO PostEsign] error make client do")
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	return nil

}
