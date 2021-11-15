package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ISniperRepository is
type ISniperRepository interface {
	LoginMehong(ctx context.Context, c models.EncryptedCookies) (response models.ResponseSniper, err error)
}

// SniperRepository is
type SniperRepository struct {
	DB    infrastructures.ISQLConnection
	Redis infrastructures.IRedis
}

//URLSniper url for sniper
var URLSniper = viper.GetString("sniper.url")

//AppNameSniper app name
var AppNameSniper = "sniper"

// Tx init a new transaction
func (r *SniperRepository) Tx() (tx *dbr.Tx, err error) {
	db := r.DB.EsignWrite()
	tx, err = db.Begin()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  nil,
		}).Error("Error begin transaction")
	}

	return
}

// LoginMehong store agent type data to database
func (r *SniperRepository) LoginMehong(ctx context.Context, c models.EncryptedCookies) (response models.ResponseSniper, err error) {

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	///login/other-mehong
	// Set Body
	body := map[string]interface{}{
		"mehong1": c.Cookies1,
		"mehong2": c.Cookies2,
		"appname": AppNameSniper,
	}
	requestBody, err := helpers.SetBody(body)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest("POST", URLSniper+"/login/other-mehong", requestBody)

	res, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"code":  5500,
			"error": err,
			"event": "error_response_from_sniper",
			"func":  "sniperRepository_LoginMehong",
		})
		return response, err
	}
	defer res.Body.Close()

	resp, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(resp, &response)

	return
}
