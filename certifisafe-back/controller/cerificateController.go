package controller

import (
	"certifisafe-back/model"
	"certifisafe-back/service"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type CertificateHandler struct {
	service service.ICertificateService
}

func NewCertificateHandler(cs service.ICertificateService) *CertificateHandler {
	return &CertificateHandler{service: cs}
}

func (ch *CertificateHandler) UpdateCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	var certificate model.Certificate
	err := json.NewDecoder(r.Body).Decode(&certificate)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	certificate, err = ch.service.UpdateCertificate(int32(id), certificate)

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	err = json.NewEncoder(w).Encode(certificate)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getErrorStatus(err error) int {
	if errors.Is(err, service.ErrIDIsNotValid) ||
		errors.Is(err, service.ErrIssuerNameIsNotValid) ||
		errors.Is(err, service.ErrFromIsNotValid) ||
		errors.Is(err, service.ErrToIsNotValid) ||
		errors.Is(err, service.ErrSubjectNameIsNotValid) ||
		errors.Is(err, service.ErrSubjectPublicKeyIsNotValid) ||
		errors.Is(err, service.ErrIssuerIdIsNotValid) ||
		errors.Is(err, service.ErrSubjectIdIsNotValid) ||
		errors.Is(err, service.ErrSignatureIsNotValid) {

		return http.StatusBadRequest
	} else if errors.Is(err, service.ErrCertificateNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (ch *CertificateHandler) CreateCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var certificate model.Certificate
	err := json.NewDecoder(r.Body).Decode(&certificate)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	certificate, err = ch.service.CreateCertificate(certificate)

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	err = json.NewEncoder(w).Encode(certificate)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (ch *CertificateHandler) GetCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	certificate, err := ch.service.GetCertificate(int32(id))

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	err = json.NewEncoder(w).Encode(certificate)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (ch *CertificateHandler) DeleteCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	var certificate model.Certificate
	err := json.NewDecoder(r.Body).Decode(&certificate)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	err = ch.service.DeleteCertificate(int32(id))

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Successfully deleted"))
}
