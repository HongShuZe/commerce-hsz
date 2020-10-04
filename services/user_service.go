package services

import (
	"commerce-hsz/datamodels"
	"commerce-hsz/repositories"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type IUserService interface {
	PwdSuccess(userName string, pwd string) (*datamodels.User, bool)
	AddUser(user *datamodels.User) (int64, error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

// 判断用户登录密码是否正确
func (u *UserService)PwdSuccess(userName string, pwd string) (*datamodels.User, bool)  {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return &datamodels.User{}, false
	}

	isOk, _ := ValidatePassword(pwd, user.Password)
	if !isOk {
		return &datamodels.User{}, false
	}

	return user, true
}

// 添加新用户
func (u *UserService) AddUser(user *datamodels.User) (int64, error) {
	pwdBytes, err := GeneratePassword(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = string(pwdBytes)
	return u.UserRepository.Insert(user)
}

// 加密
func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

//
func ValidatePassword(password string, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return false, errors.New("密码错误")
	}
	return true, nil
}