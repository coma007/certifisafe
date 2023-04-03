package controller

import (
	"certifisafe-back/model"
	"certifisafe-back/service"
	"certifisafe-back/utils"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"math/big"
	"net/http"
)

type CertificateHandler struct {
	service service.ICertificateService
}

func NewCertificateHandler(cs service.ICertificateService) *CertificateHandler {
	return &CertificateHandler{service: cs}
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

	//var certificate x509.Certificate
	//err := json.NewDecoder(r.Body).Decode(&certificate)
	//if err != nil {
	//	http.Error(w, "error when decoding json", http.StatusInternalServerError)
	//	return
	//}
	subject := pkix.Name{
		Country:            nil,
		Organization:       nil,
		OrganizationalUnit: nil,
		PostalCode:         nil,
		CommonName:         "",
		Names:              nil,
	}
	kind := model.INTERMEDIATE

	certificate, err := ch.service.CreateCertificate(subject, big.Int{}, kind)

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(certificate)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func (ch *CertificateHandler) GetCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := utils.StringToBigInt(ps.ByName("id"))

	certificate, err := ch.service.GetCertificate(id)

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(certificate)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func (ch *CertificateHandler) GetCertificates(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	certificates, err := ch.service.GetCertificates()

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(certificates)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func (ch *CertificateHandler) DeleteCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//id, _ := utils.StringToBigInt(ps.ByName("id"))

	var certificate model.Certificate
	err := json.NewDecoder(r.Body).Decode(&certificate)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	//err = ch.service.DeleteCertificate(int64(id))

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Successfully deleted"))
}

func (ch *CertificateHandler) IsValid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := utils.StringToBigInt(ps.ByName("id"))

	result, err := ch.service.IsValid(id)
	fmt.Print(result)
	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	if result {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

func (ch *CertificateHandler) Generate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	subject := pkix.Name{
		Country:            nil,
		Organization:       nil,
		OrganizationalUnit: nil,
		PostalCode:         nil,
		CommonName:         "",
		Names:              nil,
	}
	root, err := ch.service.CreateCertificate(subject, big.Int{}, model.ROOT)
	intermidiate, err := ch.service.CreateCertificate(subject, *root.SerialNumber, model.INTERMEDIATE)
	leaf, err := ch.service.CreateCertificate(subject, *intermidiate.SerialNumber, model.END)

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode([]x509.Certificate{root, intermidiate, leaf})
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}
