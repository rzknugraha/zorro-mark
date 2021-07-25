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
	"time"

	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/sirupsen/logrus"
)

// IEsignRepository is
type IEsignRepository interface {
	PostEsign(ctx context.Context, values map[string]io.Reader) (err error)
}

// EsignRepository is
type EsignRepository struct {
	DB infrastructures.ISQLConnection
}

// PostEsign post to bsre
func (r *EsignRepository) PostEsign(ctx context.Context, values map[string]io.Reader) (err error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  values,
				}).Error("[REPO PostEsign] error create form file")
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  values,
				}).Error("[REPO PostEsign] error create form file another field")
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  values,
			}).Error("[REPO PostEsign] error copy form file")
			return err
		}

	}

	w.Close()
	req, err := http.NewRequest("POST", "http://192.168.1.31/api/sign/pdf", &b)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  values,
		}).Error("[REPO PostEsign] error post data")
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.SetBasicAuth("admin", "qwerty")
	rsp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  values,
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
