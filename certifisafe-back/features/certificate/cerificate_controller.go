package certificate

import (
	"bufio"
	"certifisafe-back/utils"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type CertificateController struct {
	service CertificateService
}

func NewCertificateController(cs CertificateService) *CertificateController {
	return &CertificateController{service: cs}
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

	utils.ReturnResponse(w, err, certificates, http.StatusOK)
}

func (ch *CertificateController) DeleteCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadCertificateIDFromUrl(w, ps)
	if err != nil {
		return
	}

	err = ch.service.DeleteCertificate(id.Uint64())
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, nil, http.StatusOK)
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
