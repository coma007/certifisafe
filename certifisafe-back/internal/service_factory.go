package internal

import (
	"certifisafe-back/features/auth"
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/request"
)

type ServiceFactory interface {
	InitServices()
}

type DefaultServiceFactory struct {
	repoFactory        *DefaultRepositoryFactory
	AuthService        auth.AuthService
	OAuthService       auth.OauthService
	CertificateService certificate.CertificateService
	RequestService     request.RequestService
}

func NewDefaultServiceFactory(repoFactory DefaultRepositoryFactory) *DefaultServiceFactory {
	return &DefaultServiceFactory{
		repoFactory: &repoFactory,
	}
}

func (serviceFactory *DefaultServiceFactory) InitServices() {
	serviceFactory.AuthService = auth.NewDefaultAuthService(serviceFactory.repoFactory.UserRepository, serviceFactory.repoFactory.PasswordRecoveryRepository, serviceFactory.repoFactory.VerificationRepository)
	serviceFactory.OAuthService = auth.NewDefaultOauthService(serviceFactory.AuthService, serviceFactory.repoFactory.UserRepository)
	serviceFactory.CertificateService = certificate.NewDefaultCertificateService(serviceFactory.repoFactory.CertificateDBRepository, serviceFactory.repoFactory.CertificateFileStoreRepository, serviceFactory.repoFactory.UserRepository)
	serviceFactory.RequestService = request.NewDefaultRequestService(serviceFactory.repoFactory.RequestRepository, serviceFactory.CertificateService, serviceFactory.repoFactory.UserRepository)
}
