package internal

import (
	"certifisafe-back/features/auth"
	"certifisafe-back/features/certificate"
	"certifisafe-back/features/password_recovery"
	"certifisafe-back/features/request"
	"certifisafe-back/features/user"
	"gorm.io/gorm"
)

type RepositoryFactory interface {
	InitRepositories()
}

type DefaultRepositoryFactory struct {
	db                             *gorm.DB
	UserRepository                 user.UserRepository
	VerificationRepository         auth.VerificationRepository
	PasswordRecoveryRepository     password_recovery.PasswordRecoveryRepository
	CertificateDBRepository        certificate.CertificateRepository
	CertificateFileStoreRepository certificate.FileStoreCertificateRepository
	RequestRepository              request.RequestRepository
}

func NewDefaultRepositoryFactory(db *gorm.DB) *DefaultRepositoryFactory {
	return &DefaultRepositoryFactory{
		db: db,
	}
}

func (repoFactory *DefaultRepositoryFactory) InitRepositories() {
	repoFactory.UserRepository = user.NewDefaultUserRepository(repoFactory.db)
	repoFactory.VerificationRepository = auth.NewDefaultVerificationRepository(repoFactory.db)
	repoFactory.PasswordRecoveryRepository = password_recovery.NewDefaultPasswordRecoveryRepository(repoFactory.db)
	repoFactory.CertificateDBRepository = certificate.NewDefaultCertificateRepository(repoFactory.db)
	repoFactory.CertificateFileStoreRepository = certificate.NewDefaultFileStoreCertificateRepository()
	repoFactory.RequestRepository = request.NewDefaultRequestRepository(repoFactory.db, repoFactory.CertificateDBRepository, repoFactory.UserRepository)
}
