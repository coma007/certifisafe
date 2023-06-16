package auth

import (
	"certifisafe-back/features/user"
	"certifisafe-back/utils"
	"encoding/json"
	"gorm.io/gorm"
)

type OauthService interface {
	AuthenticateUser(data []byte) (string, error)
}

type DefaultOauthService struct {
	authService    AuthService
	userRepository user.UserRepository
}

func NewDefaultOauthService(authService AuthService, userRepository user.UserRepository) *DefaultOauthService {
	return &DefaultOauthService{
		authService:    authService,
		userRepository: userRepository,
	}
}

func (service *DefaultOauthService) AuthenticateUser(data []byte) (string, error) {
	var oauthUserMap map[string]string
	json.Unmarshal(data, &oauthUserMap)
	oauthUser := &user.User{
		Email:     oauthUserMap["email"],
		Password:  oauthUserMap["id"],
		FirstName: oauthUserMap["given_name"],
		LastName:  oauthUserMap["family_name"],
		Phone:     "",
		IsAdmin:   false,
		IsActive:  true,
	}

	_, err := service.authService.GetUserByEmail(oauthUser.Email)
	if err == gorm.ErrRecordNotFound {
		err = service.createUser(oauthUser)
		if err != nil {
			return "", err
		}
	}
	return service.authService.Login(oauthUserMap["email"], oauthUserMap["id"])
}

func (service *DefaultOauthService) createUser(oauthUser *user.User) error {
	passwordBytes, err := service.authService.HashToken(oauthUser.Password)
	utils.CheckError(err)
	oauthUser.Password = string(passwordBytes)
	_, err = service.userRepository.CreateUser(*oauthUser)
	if err != nil {
		return err
	}
	return nil
}
