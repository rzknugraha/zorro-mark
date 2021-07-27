package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// IUserService is
type IUserService interface {
	StoreUser(models.User) error
	UpdateUser(IDuser int, data models.User) (err error)
	FindUserByIDDPR(IDDPR int) (user models.User, err error)
	Login(l models.Login) (result models.TokenResp, err error)
	Profile(ctx context.Context, NIP string) (Response *helpers.JSONResponse, err error)
	UpdateFile(ctx context.Context, file multipart.File, oldName string, IDUser int, fileTypeReq string) (Response *helpers.JSONResponse, err error)
}

// UserService is
type UserService struct {
	UserRepository repositories.IUserRepository
}

// InitUserService init
func InitUserService() *UserService {
	NewUserRepository := new(repositories.UserRepository)
	NewUserRepository.DB = &infrastructures.SQLConnection{}

	UserService := new(UserService)
	UserService.UserRepository = NewUserRepository

	return UserService
}

// StoreUser is
func (p *UserService) StoreUser(data models.User) (err error) {
	_, err = p.UserRepository.StoreUser(data)
	return err
}

// UpdateUser is
func (p *UserService) UpdateUser(IDuser int, data models.User) (err error) {
	err = p.UserRepository.UpdateUserByID(IDuser, data)
	return err
}

// FindUserByIDDPR is
func (p *UserService) FindUserByIDDPR(IDDPR int) (user models.User, err error) {

	user, err = p.UserRepository.GetUserByIDDPR(IDDPR)
	fmt.Println(user)
	return
}

// Login is
func (p *UserService) Login(l models.Login) (token models.TokenResp, err error) {

	user, err := p.UserRepository.GetUserByNIPstore(l.Nip)
	fmt.Println(user)

	if err != nil {
		return
	}

	if user.Nama == "" {
		logrus.WithFields(logrus.Fields{
			"code":  4400,
			"error": err,
			"data":  user,
		}).Error("[Service Login] Wrong Username Or Password")
		return token, errors.New("Wrong Username Or Password")
	}
	//$2a$10$LT/y2441Q.rqqjWbR./9JOVWQoL1Dc6dtNRpfy6TrTx/H6XUX/A0e
	password := []byte(l.Password)
	// hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), password)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  password,
		}).Error("[Service Login] error hashing password")
		return token, errors.New("Wrong Username Or Password")
	}
	token, err = helpers.GenerateToken(user)

	return
}

// Profile is
func (p *UserService) Profile(ctx context.Context, NIP string) (Response *helpers.JSONResponse, err error) {

	Filter := map[string]interface{}{
		"nip": NIP,
	}
	user, err := p.UserRepository.FindOneUser(ctx, Filter)
	if err != nil {
		return
	}

	if user.ID == 0 {
		return &helpers.JSONResponse{
			Code:    4400,
			Message: "Not Found",
			Data:    nil,
		}, nil
	}

	return &helpers.JSONResponse{
		Code:    2200,
		Message: "Found",
		Data:    user,
	}, nil
}

// UpdateFile is
func (p *UserService) UpdateFile(ctx context.Context, file multipart.File, oldName string, IDUser int, fileTypeReq string) (Response *helpers.JSONResponse, err error) {

	trimSpace := strings.ReplaceAll(oldName, " ", "")
	path := viper.GetString("storage.path")

	fileResp := models.UploadResp{
		FileName: fmt.Sprintf("%d-%s", time.Now().UnixNano(), trimSpace),
	}

	fullPath := fmt.Sprintf("/%s/%s/%s", path, fileTypeReq, fileResp.FileName)

	dst, err := os.Create("." + fullPath)
	if err != nil {
		return

	}

	defer dst.Close()
	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return
	}

	tx, err := p.UserRepository.Tx()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  fileResp,
		}).Error("[Service UpdateFile] error create tx")
		return
	}

	defer tx.RollbackUnlessCommitted()
	// TimeNow := time.Now()

	condition := map[string]interface{}{
		"id": IDUser,
	}
	updatePayload := map[string]interface{}{
		fileTypeReq: fullPath,
	}

	res, err := p.UserRepository.UpdateUserCond(ctx, tx, condition, updatePayload)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  updatePayload,
		}).Error("[Service UpdateFile] error update file profile")
		return
	}
	tx.Commit()

	if res == 0 {
		return &helpers.JSONResponse{
			Code:    4400,
			Message: "Failed Update File",
			Data:    fileResp,
		}, nil
	}

	return &helpers.JSONResponse{
		Code:    2200,
		Message: "Success",
		Data:    fileResp,
	}, nil
}
