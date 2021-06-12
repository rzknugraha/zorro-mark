package services

import (
	"errors"
	"fmt"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// IUserService is
type IUserService interface {
	StoreUser(models.User) error
	UpdateUser(IDuser int, data models.User) (err error)
	FindUserByIDDPR(IDDPR int) (user models.User, err error)
	Login(l models.Login) (result models.TokenResp, err error)
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
