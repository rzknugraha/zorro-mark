package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var expired = time.Duration(22) * time.Hour
var appName = viper.GetString("jwt.application_name")
var signMethod = jwt.SigningMethodHS256
var signKey = []byte(viper.GetString("jwt.signature_key"))

//GenerateToken generate token JWT
func GenerateToken(user models.User) (t models.TokenResp, err error) {

	claims := models.MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    appName,
			ExpiresAt: time.Now().Add(expired).Unix(),
		},
		Nip:        user.NIP,
		ID:         user.ID,
		Name:       user.Nama,
		Position:   user.NamaJabatan,
		IdentityNO: user.Ktp,
	}

	token := jwt.NewWithClaims(
		signMethod,
		claims,
	)

	signedToken, err := token.SignedString(signKey)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  claims,
		}).Error("[Helper GenerateToken] Error Generate token")
		return
	}

	t.Token = signedToken
	t.Expired = time.Now().Add(expired).UTC().Format(time.RFC1123)

	t.DataUser.ID = user.ID
	t.DataUser.Name = user.Nama
	t.DataUser.IdentityNO = user.Ktp
	t.DataUser.Nip = user.NIP
	t.DataUser.Position = user.NamaJabatan

	return
}
