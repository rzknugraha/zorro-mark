package services

import (
	"fmt"

	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
)

// IUserService is
type IUserService interface {
	StoreUser(models.User) error
	UpdateUser(IDuser int, data models.User) (err error)
	FindUserByIDDPR(IDDPR int) (user models.User, err error)
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
