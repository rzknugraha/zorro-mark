package wrapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/guregu/null"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

//SyncUser define struct start here
type SyncUser struct {
	UserService services.IUserService
}

//InitSyncUser init sync user
func InitSyncUser() *SyncUser {
	UserService := services.InitUserService()

	wrapper := new(SyncUser)
	wrapper.UserService = UserService
	return wrapper
}

//Run start command
func (w *SyncUser) Run() {

	logrus.Info("start sync")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var userDPR []models.UserDPR

	url := viper.GetString("url.dpruser")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  nil,
		}).Error("error creating new request user")
		return
	}
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)

	res, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  nil,
		}).Error("error creating new request get user")
		return
	}

	fmt.Println(res)

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(data, &userDPR)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  res.Body,
		}).Error("error unmarshaling response")
		return
	}
	count := 0
	for _, y := range userDPR {

		intIDDPR, _ := strconv.Atoi(y.ID)
		intIDSatker, _ := strconv.Atoi(y.IDSatker)
		intIDSubSatker, _ := strconv.Atoi(y.IDSubSatker)
		userDB := models.User{
			Nama:        y.Nama,
			IDDpr:       intIDDPR,
			Ktp:         y.KTP,
			NamaJabatan: y.NamaJabatan,
			NamaSatker:  y.NamaSatker,
			NIP:         y.Nip,
			IDSatker:    intIDSatker,
			IDSubSatker: intIDSubSatker,
		}

		getUser, err := w.UserService.FindUserByIDDPR(userDB.IDDpr)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  userDB,
			}).Error("error get user")

		}

		if getUser.ID == 0 {
			fmt.Println(getUser)
			logrus.Info("User Not Found Try Insert")
			password := []byte("0987654321")
			hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  password,
				}).Error("error hashing password")
				return
			}

			userDB.Password = string(hashedPassword)
			userDB.Status = 1
			userDB.Role = null.StringFrom("user")

			err = w.UserService.StoreUser(userDB)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  password,
				}).Error("error store user")
				return
			}
			logrus.Info("success insert " + userDB.NIP)

		} else {
			userDB.Role = null.StringFrom("user")
			fmt.Println(userDB)
			err = w.UserService.UpdateUser(getUser.ID, userDB)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  getUser,
				}).Error("error update user")
				return
			}
			logrus.Info("success update " + userDB.NIP)
		}
		stringCount := strconv.Itoa(count)
		logrus.Info("iterate  " + stringCount)
		count++
	}
	return

}
