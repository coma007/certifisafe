package internal

import (
	"certifisafe-back/features/auth"
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/request"
)

type ControllerFactory interface {
	InitControllers()
}

type DefaultControllerFactory struct {
	serviceFactory        *DefaultServiceFactory
	AuthController        auth.AuthController
	CertificateController certificate.CertificateController
	RequestController     request.RequestController
}

func NewDefaultControllerFactory(serviceFactory DefaultServiceFactory) *DefaultControllerFactory {
	return &DefaultControllerFactory{
		serviceFactory: &serviceFactory,
	}
}

func (controllerFactory *DefaultControllerFactory) InitControllers() {
	controllerFactory.AuthController = *auth.NewAuthController((controllerFactory.serviceFactory.AuthService))
	controllerFactory.CertificateController = *certificate.NewCertificateController(controllerFactory.serviceFactory.CertificateService, controllerFactory.serviceFactory.AuthService)
	controllerFactory.RequestController = *request.NewRequestController(controllerFactory.serviceFactory.RequestService, controllerFactory.serviceFactory.CertificateService, controllerFactory.serviceFactory.AuthService)
}
