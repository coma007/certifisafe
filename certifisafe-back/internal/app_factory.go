package internal

import (
	"gorm.io/gorm"
)

type AppFactory interface {
	InitApp()
}

type DefaultAppFactory struct {
	db          *gorm.DB
	Repos       *DefaultRepositoryFactory
	Services    *DefaultServiceFactory
	Controllers *DefaultControllerFactory
}

func NewDefaultAppFactory(db *gorm.DB) *DefaultAppFactory {
	return &DefaultAppFactory{
		db: db,
	}
}

func (appFactory *DefaultAppFactory) InitApp() {
	appFactory.Repos = NewDefaultRepositoryFactory(appFactory.db)
	appFactory.Repos.InitRepositories()

	appFactory.Services = NewDefaultServiceFactory(*appFactory.Repos)
	appFactory.Services.InitServices()

	appFactory.Controllers = NewDefaultControllerFactory(*appFactory.Services)
	appFactory.Controllers.InitControllers()
}
