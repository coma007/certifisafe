package controller

import (
	"certifisafe-back/dto"
	"certifisafe-back/model"
	"certifisafe-back/service"
	"certifisafe-back/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
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

	var request dto.NewRequestDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	certificate, err := ch.service.CreateCertificate(request)

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

	certificate, err := ch.service.GetCertificate(id.Uint64())

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

	// TODO: response has content, fix header
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Successfully deleted"))
}

func (ch *CertificateHandler) IsValid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := utils.StringToBigInt(ps.ByName("id"))

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

func (ch *CertificateHandler) Generate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	rootDTO := &dto.NewRequestDTO{
		ParentCertificate: nil,
		Certificate: &dto.CertificateDTO{
			Serial:      nil,
			Name:        "Root",
			ValidFrom:   time.Now(),
			ValidTo:     time.Now(),
			IssuerName:  "",
			SubjectName: "",
			Status:      "",
			Type:        "ROOT",
		},
		Datetime: time.Time{},
	}

	root, err := ch.service.CreateCertificate(*rootDTO)
	if err != nil {
		panic(err)
	}
	//rootSerial := new(big.Int)
	//rootSerial.SetString(root.Serial, 10)
	//rootCreated, err := ch.service.GetCertificate(*rootSerial)
	if err != nil {
		panic(err)
	}

	intermediateDTO := &dto.NewRequestDTO{
		ParentCertificate: &root,
		Certificate: &dto.CertificateDTO{
			Serial:      nil,
			Name:        "SUB",
			ValidFrom:   time.Now(),
			ValidTo:     time.Now(),
			IssuerName:  "",
			SubjectName: "",
			Status:      "",
			Type:        "INTERMEDIATE",
		},
		Datetime: time.Time{},
	}
	intermidiate, err := ch.service.CreateCertificate(*intermediateDTO)

	leafDTO := &dto.NewRequestDTO{
		ParentCertificate: &intermidiate,
		Certificate: &dto.CertificateDTO{
			Serial:      nil,
			Name:        "LEAF",
			ValidFrom:   time.Now(),
			ValidTo:     time.Now(),
			IssuerName:  "",
			SubjectName: "",
			Status:      "",
			Type:        "END",
		},
		Datetime: time.Time{},
	}

	leaf, err := ch.service.CreateCertificate(*leafDTO)

	if err != nil {
		http.Error(w, err.Error(), getErrorStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode([]dto.CertificateDTO{root, intermidiate, leaf})
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}
