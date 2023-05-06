package certificate

import (
	"certifisafe-back/utils"
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

	utils.ReturnResponse(w, err, CertificatesToDTOs(certificates), http.StatusOK)
}

func (ch *CertificateController) WithdrawCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadCertificateIDFromUrl(w, ps)
	if err != nil {
		return
	}
	var certificate CertificateDTO
	certificate, err = ch.service.WithdrawCertificate(id.Uint64())
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	utils.ReturnResponse(w, err, certificate, http.StatusOK)
}

func (ch *CertificateController) IsValid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO finish - Bobi
	id, err := utils.ReadCertificateIDFromUrl(w, ps)
	if err != nil {
		return
	}

	result, err := ch.service.IsValid(id.Uint64())
	fmt.Print(result)
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	// TODO: response has content, fix header
	w.WriteHeader(http.StatusNoContent)
	if result {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}
