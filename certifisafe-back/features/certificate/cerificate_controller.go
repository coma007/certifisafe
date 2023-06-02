package certificate

import (
	"bufio"
	"certifisafe-back/features/auth"
	"certifisafe-back/utils"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type CertificateController struct {
	certificateService CertificateService
	authService        auth.AuthService
}

func NewCertificateController(certificateService CertificateService, authService auth.AuthService) *CertificateController {
	return &CertificateController{certificateService: certificateService, authService: authService}
}

func getErrorStatus(err error) int {
	if errors.Is(err, ErrIDIsNotValid) ||
		errors.Is(err, ErrIssuerNameIsNotValid) ||
		errors.Is(err, ErrFromIsNotValid) ||
		errors.Is(err, ErrToIsNotValid) ||
		errors.Is(err, ErrSubjectNameIsNotValid) ||
		errors.Is(err, ErrSubjectPublicKeyIsNotValid) ||
		errors.Is(err, ErrIssuerIdIsNotValid) ||
		errors.Is(err, ErrSubjectIdIsNotValid) ||
		errors.Is(err, ErrSignatureIsNotValid) {

		return http.StatusBadRequest
	} else if errors.Is(err, ErrCertificateNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (controller *CertificateController) GetCertificate(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadCertificateIDFromUrl(w, r)
	if err != nil {
		return
	}

	certificate, err := controller.certificateService.GetCertificate(id.Uint64())
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, certificate, http.StatusOK)
}

func (controller *CertificateController) GetCertificates(w http.ResponseWriter, r *http.Request) {
	certificates, err := controller.certificateService.GetCertificates()
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, CertificatesToDTOs(certificates), http.StatusOK)
}

func (controller *CertificateController) DownloadCertificate(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadCertificateIDFromUrl(w, r)
	if err != nil {
		return
	}
	user := controller.authService.GetUserFromToken(r.Header.Get("Authorization"))
	publicPath, privatePath, err := controller.certificateService.GetCertificateFiles(id.Uint64(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !utils.AddFileToResponse(w, publicPath) {
		return
	}
	if privatePath != "" && !utils.AddFileToResponse(w, privatePath) {
		return
	}

	w.Header().Set("  Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=certificate.zip")
}

func (controller *CertificateController) WithdrawCertificate(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadCertificateIDFromUrl(w, r)
	if err != nil {
		return
	}
	user := controller.authService.GetUserFromToken(r.Header.Get("Authorization"))

	var certificate CertificateDTO
	certificate, err = controller.certificateService.WithdrawCertificate(id.Uint64(), user)
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, certificate, http.StatusOK)
}

func (controller *CertificateController) IsValid(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadCertificateIDFromUrl(w, r)
	if err != nil {
		return
	}

	result, err := controller.certificateService.IsValidById(id.Uint64())
	fmt.Print(result)
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, &result, http.StatusOK)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (controller *CertificateController) IsValidFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.ReturnResponse(w, err, nil, http.StatusBadRequest)
		return
	}

	// KEY OF THE MULTIPART FOR FILE MUST BE "file"
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.ReturnResponse(w, err, nil, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes := make([]byte, handler.Size)
	buf := bufio.NewReader(file)
	if strings.Split(handler.Filename, ".")[1] != "crt" ||
		!stringInSlice(handler.Header.Get("Content-Type"), []string{"application/pkix-cert", "application/x-x509-ca-cert"}) {
		http.Error(w, "file type not allowed", http.StatusBadRequest)
		return
	}

	_, err = buf.Read(fileBytes)
	if err != nil {
		utils.ReturnResponse(w, err, nil, http.StatusBadRequest)
		return
	}
	block, _ := pem.Decode(fileBytes)
	if block == nil {
		err = errors.New("File uploaded is not a certificate")
		utils.ReturnResponse(w, err, nil, http.StatusBadRequest)
		return
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		utils.ReturnResponse(w, err, nil, http.StatusBadRequest)
		return
	}
	result, err := controller.certificateService.IsValid(*cert)
	if err != nil {
		utils.ReturnResponse(w, err, nil, http.StatusBadRequest)
		return
	}

	utils.ReturnResponse(w, err, &result, http.StatusOK)
}
