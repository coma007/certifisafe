package certificate

import (
	"bufio"
	"certifisafe-back/features/auth"
	"certifisafe-back/utils"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type CertificateController struct {
	service     CertificateService
	authService auth.AuthService
}

func NewCertificateController(cs CertificateService, as auth.AuthService) *CertificateController {
	return &CertificateController{service: cs, authService: as}
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

func (ch *CertificateController) GetCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadCertificateIDFromUrl(w, ps)
	if err != nil {
		return
	}

	certificate, err := ch.service.GetCertificate(id.Uint64())
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, certificate, http.StatusOK)
}

func (ch *CertificateController) GetCertificates(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	certificates, err := ch.service.GetCertificates()
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, CertificatesToDTOs(certificates), http.StatusOK)
}

func (ch *CertificateController) DownloadCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadCertificateIDFromUrl(w, ps)
	if err != nil {
		return
	}
	user := ch.authService.GetUserFromToken(r.Header.Get("Authorization"))

	public, private, err := ch.service.GetCertificateFiles(id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	//w.Header().Set("Content-Disposition", "attachment; filename="+id+".pem")

	// Write file contents to response body
	//io.Copy(w, file)
}

func (ch *CertificateController) WithdrawCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadCertificateIDFromUrl(w, ps)
	if err != nil {
		return
	}
	user := ch.authService.GetUserFromToken(r.Header.Get("Authorization"))

	var certificate CertificateDTO
	certificate, err = ch.service.WithdrawCertificate(id.Uint64(), user)
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, certificate, http.StatusOK)
}

func (ch *CertificateController) IsValid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadCertificateIDFromUrl(w, ps)
	if err != nil {
		return
	}

	result, err := ch.service.IsValidById(id.Uint64())
	fmt.Print(result)
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, &result, http.StatusOK)
}

func (ch *CertificateController) IsValidFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	_, err = bufio.NewReader(file).Read(fileBytes)
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
	result, err := ch.service.IsValid(*cert)
	if err != nil {
		utils.ReturnResponse(w, err, nil, http.StatusBadRequest)
		return
	}

	utils.ReturnResponse(w, err, &result, http.StatusOK)
}
